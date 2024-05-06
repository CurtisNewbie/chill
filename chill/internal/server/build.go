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
	StatusBuilding   = "BUILDING"
	StatusFailed     = "FAILED"
)

var (
	_builds     Builds
	_buildsOnce sync.Once
)

var (
	buildPool      = miso.NewAsyncPool(60, 30)
	buildStatusMap = sync.Map{}
)

func init() {
	miso.SetDefProp(PropScriptsBaseFolder, "./")
}

type Builds struct {
	Builds []BuildConf `mapstructure:"build"`
}

func (b Builds) Find(name string) (BuildConf, bool) {
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

func LoadBuilds() Builds {
	_buildsOnce.Do(func() {
		var builds Builds
		miso.UnmarshalFromProp(&builds)
		_builds = builds
	})
	return _builds
}

func InitBuildStatusMap(bu Builds) {
	for _, b := range bu.Builds {
		buildStatusMap.Store(b.Name, false)
	}
}

func CheckBuildsConf(bu Builds) error {
	names := miso.NewSet[string]()
	for _, b := range bu.Builds {
		if b.Name == "" {
			return miso.NewErrf("build name should be empty")
		}
		if !names.Add(b.Name) {
			return miso.NewErrf("build name contains duplicate, '%s' already exists", b.Name)
		}
		{
			actual := len([]rune(b.Name))
			if actual > 128 {
				return miso.NewErrf("build name exceeds 128 characters (%d characters), '%s'", actual, b.Name)
			}
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
	builds := LoadBuilds()
	return miso.NewPageQuery[ApiListBuildInfoRes]().
		WithPage(page).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Table("build_info").Order("id desc")
		}).
		ForEach(func(t ApiListBuildInfoRes) ApiListBuildInfoRes {
			b, ok := builds.Find(t.Name)
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
				t.Triggerable = true

				if v, ok := buildStatusMap.Load(b.Name); ok && v.(bool) {
					t.Status = StatusBuilding
					t.Triggerable = false
				}
			}
			return t
		}).
		Exec(rail, db)
}

func InitBuildInfo(rail miso.Rail, builds Builds, db *gorm.DB) error {
	for _, b := range builds.Builds {
		t := db.Exec(`INSERT IGNORE INTO build_info (name, status) VALUES (?,?)`, b.Name, StatusSuccessful)
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
			tx = tx.Table("build_log").
				Order("id desc")
			if req.Name != "" {
				tx = tx.Where("name = ?", req.Name)
			}
			return tx
		}).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "name", "build_no", "status", "ctime")
		}).
		Exec(rail, db)
}

func TriggerBuild(rail miso.Rail, req ApiTriggerBuildReq, db *gorm.DB) error {
	b, ok := LoadBuilds().Find(req.Name)
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
		err := tx.Exec(`UPDATE build_info SET status = ?, utime = ? WHERE name = ?`, status, miso.Now(), name).Error
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

func QryBuildHistDetails(rail miso.Rail, db *gorm.DB, req ApiQryBuildHistDetailReq) (ApiQryBuildHistDetailRes, error) {

	var his ApiListBuildHistoryRes
	err := db.Raw(`SELECT * FROM build_log WHERE build_no = ?`, req.BuildNo).Scan(&his).Error
	if err != nil {
		return ApiQryBuildHistDetailRes{}, fmt.Errorf("failed to query build_log, %v, %w", req.BuildNo, err)
	}

	var cl []ApiCmdLogRes
	err = db.Raw(`SELECT id, command, remark, status FROM command_log WHERE build_no = ? ORDER BY id DESC`, req.BuildNo).Scan(&cl).Error
	if err != nil {
		return ApiQryBuildHistDetailRes{}, fmt.Errorf("failed to query command_log, %v, %w", req.BuildNo, err)
	}
	if cl == nil {
		cl = []ApiCmdLogRes{}
	}
	return ApiQryBuildHistDetailRes{
		Id:          his.Id,
		Name:        his.Name,
		BuildNo:     his.BuildNo,
		Status:      his.Status,
		Ctime:       his.Ctime,
		CommandLogs: cl,
	}, nil
}
