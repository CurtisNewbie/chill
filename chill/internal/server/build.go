package server

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
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
	buildPool      = util.NewAsyncPool(60, 30)
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
	Name    string
	GitRepo string `mapstructure:"git-repo"`
	Steps   []BuildCmd
}

type BuildCmd struct {
	Script  string
	Command string
}

func (b BuildCmd) IsEmpty() bool {
	return b.Script == "" && b.Command == ""
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
	names := util.NewSet[string]()
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
	return util.ReadFileAll(filepath.Join(miso.GetPropStr(PropScriptsBaseFolder), path))
}

func ListBuildInfos(rail miso.Rail, page miso.Paging, db *gorm.DB) (miso.PageRes[ApiListBuildInfoRes], error) {
	builds := LoadBuilds()
	return miso.NewPageQuery[ApiListBuildInfoRes]().
		WithPage(page).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Table("build_info").Order("utime desc")
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
			return tx.Select("id", "name", "build_no", "status", "build_start_time start_time", "build_end_time end_time", "commit_id")
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
		stime := util.Now()
		rail = rail.NextSpan()
		buildNo := util.GenIdP("build_")
		if !buildStatusMap.CompareAndSwap(b.Name, false, true) {
			rail.Infof("Build '%s' is running, ignored", b.Name)
			return
		}
		rail.Infof("Running build '%s', buildNo: %s", b.Name, buildNo)
		defer buildStatusMap.Store(b.Name, false)
		defer miso.TimeOp(rail, time.Now(), fmt.Sprintf("build '%s'", b.Name))

		// get commit id
		var commitId string
		if b.GitRepo != "" {
			commitIdCmd := fmt.Sprintf("cd %s && git rev-parse HEAD", b.GitRepo)
			cid, err := RunBuildCmd(rail, BuildCmd{Command: commitIdCmd})
			if err == nil {
				commitId = cid
				rail.Infof("Executed %v %#v, commit_id: %s", b.Name, commitIdCmd, commitId)
			} else {
				rail.Errorf("Failed to get build commit_id, %v, '%s', %v", b.Name, commitIdCmd, err)
			}
		}

		var remark string
		var status string = StatusSuccessful
		for _, s := range b.Steps {
			out, sterr := RunBuildCmd(rail, s)
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
				remark = LastNStr(remark, 1000)
			}
			if scerr := SaveCmdLog(rail, db, buildNo, cmd, status, remark); scerr != nil {
				rail.Errorf("Failed to save command log, build: %s, %v", b.Name, scerr)
				return
			}
			if sterr != nil {
				break
			}
		}

		ubsp := UpdateBuildStatusParam{BuildNo: buildNo,
			Name:      b.Name,
			Status:    status,
			Remark:    remark,
			StartTime: stime,
			EndTime:   util.Now(),
			CommitId:  commitId,
		}
		if er := UpdateBuildStatus(rail, db, ubsp); er != nil {
			rail.Errorf("Failed to save build log, build: %s, %#v, %v", b.Name, ubsp, er)
		}
	})
	return nil
}

func RunBuildCmd(rail miso.Rail, cmd BuildCmd) (string, error) {
	if cmd.Command != "" {
		return BashRun(rail, util.UnsafeStr2Byt(cmd.Command))
	} else {
		file, err := LookupBuildScript(cmd.Script)
		if err != nil {
			return "", fmt.Errorf("failed to load build script, %v", file)
		}
		return BashRun(rail, file)
	}
}

func SaveCmdLog(rail miso.Rail, db *gorm.DB, buildNo string, cmd string, status string, remark string) error {
	return db.Exec(`INSERT INTO command_log (build_no,command,remark,status) VALUES
	(?,?,?,?)`, buildNo, cmd, remark, status).Error
}

type UpdateBuildStatusParam struct {
	BuildNo   string
	Name      string
	Status    string
	Remark    string
	CommitId  string
	StartTime util.ETime
	EndTime   util.ETime
}

func UpdateBuildStatus(rail miso.Rail, db *gorm.DB, p UpdateBuildStatusParam) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`UPDATE build_info SET status = ?, utime = ?, commit_id = ? WHERE name = ?`, p.Status, util.Now(), p.CommitId, p.Name).Error
		if err != nil {
			return fmt.Errorf("failed to update build_info, %w", err)
		}

		err = tx.Exec(`INSERT INTO build_log (build_no, name, status, remark, build_start_time, build_end_time, commit_id) VALUES (?,?,?,?,?,?,?)`,
			p.BuildNo, p.Name, p.Status, p.Remark, p.StartTime, p.EndTime, p.CommitId).Error
		if err != nil {
			return fmt.Errorf("failed to save build_log, %w", err)
		}
		return nil
	})
}

func QryBuildHistDetails(rail miso.Rail, db *gorm.DB, req ApiQryBuildHistDetailReq) (ApiQryBuildHistDetailRes, error) {

	var his ApiListBuildHistoryRes
	err := db.
		Raw(`
		SELECT id, name, build_no, status, build_start_time start_time, build_end_time end_time, commit_id
		FROM build_log WHERE build_no = ?`, req.BuildNo).
		Scan(&his).Error
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
		CommitId:    his.CommitId,
		Status:      his.Status,
		StartTime:   his.StartTime,
		EndTime:     his.EndTime,
		CommandLogs: cl,
	}, nil
}

func LastNStr(s string, n int) string {
	ru := []rune(s)
	if len(ru) <= n {
		return s
	}
	return string(ru[len(ru)-n:])
}
