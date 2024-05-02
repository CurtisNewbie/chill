package server

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

const (
	PropScriptsBaseFolder = "scripts.base-folder"

	StatusSuccessful = "SUCCESSFUL"
	StatusFailed     = "FAILED"
)

var (
	buildsConf BuildsConf
	bcOnce     sync.Once
)

var (
	buildPool      = miso.NewAsyncPool(60, 30)
	buildStatusMap = sync.Map{}
)

func init() {
	miso.SetDefProp(PropScriptsBaseFolder, "./")
}

type BuildsConf struct {
	Builds []BuildConf `mapstructure:"build"`
}

func (b BuildsConf) Find(name string) (BuildConf, bool) {
	for i, ib := range b.Builds {
		if ib.Name == name {
			return b.Builds[i], true
		}
	}
	return BuildConf{}, false
}

type BuildConf struct {
	Name  string
	Steps []BuildStep
}

type BuildStep struct {
	Script  string
	Command string
}

func LoadBuildsConf() BuildsConf {
	bcOnce.Do(func() {
		var bc BuildsConf
		miso.UnmarshalFromProp(&bc)
		buildsConf = bc
	})
	return buildsConf
}

func InitBuildStatusMap(bc BuildsConf) {
	for _, b := range bc.Builds {
		buildStatusMap.Store(b.Name, false)
	}
}

func CheckBuildsConf(bc BuildsConf) error {
	names := miso.NewSet[string]()
	for _, b := range bc.Builds {
		if b.Name == "" {
			return miso.NewErrf("build name should be empty")
		}
		if !names.Add(b.Name) {
			return miso.NewErrf("build name contains duplicate, '%s' already exists", b.Name)
		}
		if len(b.Steps) < 1 {
			return miso.NewErrf("build '%s' contains zero step, configuration illegal", b.Name)
		}
		for _, s := range b.Steps {
			if s.Command != "" && s.Script != "" {
				return miso.NewErrf("build step illegal, either configure script or command but not both")
			}
		}
	}
	return nil
}

// Lookup build script under the specified folder.
func LookupBuildScript(path string) ([]byte, error) {
	return miso.ReadFileAll(filepath.Join(miso.GetPropStr(PropScriptsBaseFolder), path))
}

func ListBuildInfos(rail miso.Rail, page miso.Paging, db *gorm.DB) (miso.PageRes[ApiListBuildInfoRes], error) {
	bc := LoadBuildsConf()
	return miso.NewPageQuery[ApiListBuildInfoRes]().
		WithPage(page).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Table("build_info").Order("id desc")
		}).
		ForEach(func(t ApiListBuildInfoRes) ApiListBuildInfoRes {
			b, ok := bc.Find(t.Name)
			if ok {
				t.BuildSteps = make([]string, 0, len(b.Steps))
				for _, st := range b.Steps {
					var ss string
					if st.Script != "" {
						ss = "file: " + st.Script
					} else {
						ss = "bash: " + st.Command
					}
					t.BuildSteps = append(t.BuildSteps, ss)
				}
			}
			return t
		}).
		Exec(rail, db)
}

func InitBuildInfo(rail miso.Rail, bc BuildsConf, db *gorm.DB) error {
	for _, b := range bc.Builds {
		t := db.Exec(`INSERT IGNORE INTO build_info (name, status) VALUES (?,?)`, b.Name, "SUCCESSFUL")
		if t.Error != nil {
			return fmt.Errorf("failed to init build_info record, name: %v, %w", b.Name, t.Error)
		}
	}
	return nil
}

func ListBuildHistory(rail miso.Rail, req ApiListBuildHistoryReq, db *gorm.DB) (miso.PageRes[ApiListBuildHistoryRes], error) {
	return miso.NewPageQuery[ApiListBuildHistoryRes]().
		WithPage(req.Paging).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Table("build_log").
				Where("name = ?", req.Name).
				Order("id desc")
		}).
		Exec(rail, db)
}

func TriggerBuild(rail miso.Rail, req ApiTriggerBuildReq, db *gorm.DB) error {
	b, ok := LoadBuildsConf().Find(req.Name)
	if !ok {
		return miso.NewErrf("Build name not found")
	}

	v, _ := buildStatusMap.Load(b.Name)
	if vb := v.(bool); vb {
		return miso.NewErrf("Build '%s' is running", b.Name)
	}

	buildPool.Go(func() {
		rail = rail.NextSpan()
		buildNo := miso.GenIdP("build_")
		if !buildStatusMap.CompareAndSwap(b.Name, false, true) {
			rail.Infof("Build '%s' is running, ignored", b.Name)
			return
		}
		rail.Infof("Running build '%s', buildNo: %s", b.Name, buildNo)
		defer buildStatusMap.Store(b.Name, false)
		defer miso.TimeOp(rail, time.Now(), fmt.Sprintf("build '%s'", b.Name))

		time.Sleep(500 * time.Millisecond)

		var remark string
		var status string = StatusSuccessful
		for _, s := range b.Steps {
			out, sterr := RunStep(rail, s)
			if sterr != nil {
				rail.Warnf("Failed to run step %#v, build: %s, %v, %v", s, b.Name, out, sterr)
			}

			var cmd string
			if s.Command != "" {
				cmd = s.Command
			} else {
				cmd = "file: " + s.Script
			}

			if sterr != nil {
				remark = sterr.Error()
				status = StatusFailed
			} else {
				remark = out
			}

			if remark != "" {
				remark = miso.MaxLenStr(remark, 1000)
			}
			if scerr := SaveCmdLog(rail, db, buildNo, b.Name, cmd, status, remark); scerr != nil {
				rail.Errorf("Failed to save command log, build: %s, %v", b.Name, scerr)
				return
			}
			if sterr != nil {
				break
			}
		}

		if er := UpdateBuildStatus(rail, db, buildNo, b.Name, status, remark); er != nil {
			rail.Errorf("Failed to save build log, build: %s, %v", b.Name, er)
		}
	})
	return nil
}

func RunStep(rail miso.Rail, step BuildStep) (string, error) {
	if step.Command != "" {
		return BashRun(rail, miso.UnsafeStr2Byt(step.Command))
	} else {
		file, err := LookupBuildScript(step.Script)
		if err != nil {
			return "", fmt.Errorf("failed to load build script, %v", file)
		}
		return BashRun(rail, file)
	}
}

func SaveCmdLog(rail miso.Rail, db *gorm.DB, buildNo string, name string, cmd string, status string, remark string) error {
	return db.Exec(`INSERT INTO command_log (build_no, build_name,command,remark,status) VALUES
	(?,?,?,?,?)`, buildNo, name, cmd, remark, status).Error
}

func UpdateBuildStatus(rail miso.Rail, db *gorm.DB, buildNo string, name string, status string, remark string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`UPDATE build_info SET status = ? WHERE name = ?`, status, name).Error
		if err != nil {
			return fmt.Errorf("failed to update build_info, %w", err)
		}

		err = tx.Exec(`INSERT INTO build_log (build_no, name, status, remark) VALUES (?,?,?,?)`, buildNo, name, status, remark).Error
		if err != nil {
			return fmt.Errorf("failed to save build_log, %w", err)
		}
		return nil
	})
}

func QryBuildHistDetails(rail miso.Rail, db *gorm.DB, req ApiQryBuildHistReq) (ApiQryBuildHistRes, error) {

	var his ApiListBuildHistoryRes
	err := db.Raw(`SELECT * FROM build_log WHERE build_no = ?`, req.BuildNo).Scan(&his).Error
	if err != nil {
		return ApiQryBuildHistRes{}, fmt.Errorf("failed to query build_log, %v, %w", req.BuildNo, err)
	}

	var cl []ApiCmdLogRes
	err = db.Raw(`SELECT id, command, remark, status FROM command_log where build_no = ?`, req.BuildNo).Scan(&cl).Error
	if err != nil {
		return ApiQryBuildHistRes{}, fmt.Errorf("failed to query command_log, %v, %w", req.BuildNo, err)
	}
	if cl == nil {
		cl = []ApiCmdLogRes{}
	}
	return ApiQryBuildHistRes{
		Id:          his.Id,
		BuildNo:     his.BuildNo,
		Status:      his.Status,
		Ctime:       his.Ctime,
		Remark:      his.Remark,
		CommandLogs: cl,
	}, nil
}
