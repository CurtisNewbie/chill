package server

import "github.com/curtisnewbie/miso/miso"

func RegisterEndpoints(rail miso.Rail) {
	miso.Get("/api/build/list-names", ApiListBuildNames)
	miso.IPost("/api/build/info/list", ApiListBuildInfos)
	miso.IPost("/api/build/trigger", ApiTriggerBuild)
	miso.IPost("/api/build/history/list", ApiListBuildHistory)
	miso.IPost("/api/build/history/detail", ApiQryBuildHistoryDetails)
}

func ApiListBuildNames(inb *miso.Inbound) ([]string, error) {
	bc := LoadBuildsConf()
	names := make([]string, 0, len(bc.Builds))
	for i := range bc.Builds {
		names = append(names, bc.Builds[i].Name)
	}
	return names, nil
}

type ApiListBuildInfoRes struct {
	Id         int
	Name       string
	Status     string
	Ctime      miso.ETime
	Utime      miso.ETime
	BuildSteps []string `gorm:"-"`
}

func ApiListBuildInfos(inb *miso.Inbound, req miso.Paging) (miso.PageRes[ApiListBuildInfoRes], error) {
	return ListBuildInfos(inb.Rail(), req, miso.GetMySQL())
}

type ApiTriggerBuildReq struct {
	Name string
}

func ApiTriggerBuild(inb *miso.Inbound, req ApiTriggerBuildReq) (any, error) {
	return nil, TriggerBuild(inb.Rail(), req, miso.GetMySQL())
}

type ApiListBuildHistoryReq struct {
	Name   string
	Paging miso.Paging
}

type ApiListBuildHistoryRes struct {
	Id      int
	BuildNo string
	Status  string
	Ctime   miso.ETime
	Remark  string
}

func ApiListBuildHistory(inb *miso.Inbound, req ApiListBuildHistoryReq) (miso.PageRes[ApiListBuildHistoryRes], error) {
	return ListBuildHistory(inb.Rail(), req, miso.GetMySQL())
}

type ApiQryBuildHistReq struct {
	BuildNo string
}

type ApiQryBuildHistRes struct {
	Id          int
	BuildNo     string
	Status      string
	Ctime       miso.ETime
	Remark      string
	CommandLogs []ApiCmdLogRes
}

type ApiCmdLogRes struct {
	Id      int
	Command string
	Remark  string
	Status  string
}

func ApiQryBuildHistoryDetails(inb *miso.Inbound, req ApiQryBuildHistReq) (ApiQryBuildHistRes, error) {
	return QryBuildHistDetails(inb.Rail(), miso.GetMySQL(), req)
}
