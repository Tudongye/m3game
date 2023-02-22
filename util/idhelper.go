package util

import (
	"fmt"
	"regexp"
	"strconv"
)

/*
AppID : EnvID.WorldID.FuncID.InsID
SvcID : EnvID.WorldID.FuncID
*/
var (
	regexAppID, _ = regexp.Compile("^([^\\.]+)\\.([^\\.]+)\\.([^\\.]+)\\.([^\\.]+)$")
	regexSvcID, _ = regexp.Compile("^([^\\.]+)\\.([^\\.]+)\\.([^\\.]+)$")
	regexAddr, _  = regexp.Compile("^([^:]+):([0-9]+)$")
)

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

func Addr2IPPort(s string) (ip string, port int, err error) {
	err = nil
	groups := regexAddr.FindStringSubmatch(s)
	if len(groups) == 0 {
		err = fmt.Errorf("Addr Parse fail %s", s)
		return
	}
	ip = groups[1]
	port, _ = strconv.Atoi(groups[2])
	return

}
