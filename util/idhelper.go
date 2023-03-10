package util

import (
	"fmt"
	"m3game/proto/pb"
	"regexp"
	"strconv"
)

/*
AppID : EnvID.WorldID.FuncID.InsID
SvcID : EnvID.WorldID.FuncID
*/
var (
	regexAppID   *regexp.Regexp
	regexSvcID   *regexp.Regexp
	regexWorldID *regexp.Regexp
	regexEnvID   *regexp.Regexp
	regexAddr    *regexp.Regexp
)

func init() {
	var err error
	if regexAppID, err = regexp.Compile("^([^\\.]+)\\.([^\\.]+)\\.([^\\.]+)\\.([^\\.]+)$"); err != nil {
		panic(fmt.Sprintf("regexAppID.Compile err %s", err))
	}
	if regexSvcID, err = regexp.Compile("^([^\\.]+)\\.([^\\.]+)\\.([^\\.]+)$"); err != nil {
		panic(fmt.Sprintf("regexSvcID.Compile err %s", err))
	}
	if regexWorldID, err = regexp.Compile("^([^\\.]+)\\.([^\\.]+)$"); err != nil {
		panic(fmt.Sprintf("regexWorldID.Compile err %s", err))
	}
	if regexEnvID, err = regexp.Compile("^([^\\.]+)$"); err != nil {
		panic(fmt.Sprintf("regexEnvID.Compile err %s", err))
	}
	if regexAddr, err = regexp.Compile("^([^:]+):([0-9]+)$"); err != nil {
		panic(fmt.Sprintf("regexAddr.Compile err %s", err))
	}
}

func AppID2Str(envid string, worldid string, funcid string, insid string) string {
	return fmt.Sprintf("%s.%s.%s.%s", envid, worldid, funcid, insid)
}

func AppStr2ID(s string) (envid string, worldid string, funcid string, insid string, err error) {
	err = nil
	groups := regexAppID.FindStringSubmatch(s)
	if len(groups) == 0 {
		err = fmt.Errorf("IDStr Parse fail %s", s)
		return
	}
	envid = groups[1]
	worldid = groups[2]
	funcid = groups[3]
	insid = groups[4]
	return
}

func SvcID2Str(envid string, worldid string, funcid string) string {
	return fmt.Sprintf("%s.%s.%s", envid, worldid, funcid)
}

func SvcStr2ID(s string) (envid string, worldid string, funcid string, err error) {
	err = nil
	groups := regexSvcID.FindStringSubmatch(s)
	if len(groups) == 0 {
		err = fmt.Errorf("IDStr Parse fail %s", s)
		return
	}
	envid = groups[1]
	worldid = groups[2]
	funcid = groups[3]
	return
}

func WorldID2Str(envid string, worldid string) string {
	return fmt.Sprintf("%s.%s", envid, worldid)
}

func WorldStr2ID(s string) (envid string, worldid string, err error) {
	err = nil
	groups := regexWorldID.FindStringSubmatch(s)
	if len(groups) == 0 {
		err = fmt.Errorf("IDStr Parse fail %s", s)
		return
	}
	envid = groups[1]
	worldid = groups[2]
	return
}

func EnvID2Str(envid string) string {
	return fmt.Sprintf("%s", envid)
}

func EnvStr2ID(s string) (envid string, err error) {
	err = nil
	groups := regexEnvID.FindStringSubmatch(s)
	if len(groups) == 0 {
		err = fmt.Errorf("IDStr Parse fail %s", s)
		return
	}
	envid = groups[1]
	return
}

func Addr2IPPort(s string) (ip string, port int, err error) {
	err = nil
	groups := regexAddr.FindStringSubmatch(s)
	if len(groups) == 0 {
		err = fmt.Errorf("Addr Parse fail %s", s)
		return
	}
	ip = groups[1]
	port, err = strconv.Atoi(groups[2])
	return

}

func RouteIns2Svc(ins *pb.RouteIns, funcid string) *pb.RouteSvc {
	return &pb.RouteSvc{
		EnvID:   ins.EnvID,
		WorldID: ins.WorldID,
		FuncID:  funcid,
		IDStr:   SvcID2Str(ins.EnvID, ins.WorldID, funcid),
	}
}

func GenSessionID(s string) string {
	return fmt.Sprintf("session-%s", s)
}
