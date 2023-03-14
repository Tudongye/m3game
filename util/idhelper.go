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
	regexAddr *regexp.Regexp
)

func init() {
	var err error
	if regexAddr, err = regexp.Compile("^([^:]+):([0-9]+)$"); err != nil {
		panic(fmt.Sprintf("regexAddr.Compile err %s", err))
	}
}
func GenSessionID(s string) string {
	return fmt.Sprintf("session-%s", s)
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
