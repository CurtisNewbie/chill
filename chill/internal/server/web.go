package server

import "github.com/curtisnewbie/miso/miso"

func RegisterEndpoints(rail miso.Rail) {
	miso.PrepareWebStaticFs(staticFs, staticFsPre)
	miso.Get("/api/build/list-names", ApiListBuildNames)
	miso.IPost("/api/build/info/list", ApiListBuildInfos)
	miso.IPost("/api/build/trigger", ApiTriggerBuild)
	miso.IPost("/api/build/history/list", ApiListBuildHistory)
	miso.IPost("/api/build/history/detail", ApiQryBuildHistoryDetails)
}

func ApiListBuildNames(inb *miso.Inbound) ([]string, error) {
	builds := LoadBuilds()
	names := make([]string, 0, len(builds.Builds))
	for i := range builds.Builds {
		names = append(names, builds.Builds[i].Name)
	}
	return names, nil
}

type ApiListBuildInfoRes struct {
	Id          int        `desc:"build info id"`
	Name        string     `desc:"build name"`
	Status      string     `desc:"last build status"`
	Ctime       miso.ETime `desc:"create time"`
	Utime       miso.ETime `desc:"update time"`
	BuildSteps  []string   `gorm:"-" desc:"build steps"`
	Triggerable bool       `gorm:"-" desc:"whether the build is triggerable"`
}

func ApiListBuildInfos(inb *miso.Inbound, req miso.Paging) (miso.PageRes[ApiListBuildInfoRes], error) {
	return ListBuildInfos(inb.Rail(), req, miso.GetMySQL())
}

type ApiTriggerBuildReq struct {
	Name string `desc:"build name" vaild:"notEmpty"`
}

func ApiTriggerBuild(inb *miso.Inbound, req ApiTriggerBuildReq) (any, error) {
	return nil, TriggerBuild(inb.Rail(), req, miso.GetMySQL())
}

type ApiListBuildHistoryReq struct {
	Name   string `desc:"build name" vaild:"notEmpty"`
	Paging miso.Paging
}

type ApiListBuildHistoryRes struct {
	Id        int        `desc:"build history id"`
	Name      string     `desc:"build name"`
	BuildNo   string     `desc:"build no"`
	Status    string     `desc:"built status"`
	StartTime miso.ETime `desc:"build start time"`
	EndTime   miso.ETime `desc:"build end time"`
}

func ApiListBuildHistory(inb *miso.Inbound, req ApiListBuildHistoryReq) (miso.PageRes[ApiListBuildHistoryRes], error) {
	return ListBuildHistory(inb.Rail(), req, miso.GetMySQL())
}

type ApiQryBuildHistDetailReq struct {
	BuildNo string `desc:"build no" vaild:"notEmpty"`
}

type ApiQryBuildHistDetailRes struct {
	Id          int            `desc:"build history id"`
	Name        string         `desc:"build name"`
	BuildNo     string         `desc:"build no"`
	Status      string         `desc:"built status"`
	StartTime   miso.ETime     `desc:"build start time"`
	EndTime     miso.ETime     `desc:"build end time"`
	Remark      string         `desc:"remark"`
	CommandLogs []ApiCmdLogRes `desc:"commands execution log"`
}

type ApiCmdLogRes struct {
	Id      int    `desc:"command log id"`
	Command string `desc:"command"`
	Remark  string `desc:"remark"`
	Status  string `desc:"execution status"`
}

func ApiQryBuildHistoryDetails(inb *miso.Inbound, req ApiQryBuildHistDetailReq) (ApiQryBuildHistDetailRes, error) {
	return QryBuildHistDetails(inb.Rail(), miso.GetMySQL(), req)
}
