package agent

import (
	"io"
)

const (
	SessionKey_Uid = "uid"
)

type Agent interface {
}

type AuthFunc func(para AuthPara) (string, error)
type CallFunc func(method string, uid string, reqbuf io.ReadCloser) (io.ReadCloser, error)

var (
	_authfunc AuthFunc // 鉴权接口，成功返回UID(玩家身份标识)
	_callfunc CallFunc // 业务接口
)

type UType int

const (
	UType_Test UType = iota + 1 // 测试类型
)

type AuthPara struct {
	Type  UType  // 玩家类型
	Token string // 鉴权信息
	Opts  map[string]string
}

type AuthRsp struct {
	Uid     string
	Session string
	EnvID   string
	WorldID string
	Opts    map[string]string
}

func SetAuther(f AuthFunc) {
	_authfunc = f
}

func Auther() AuthFunc {
	if _authfunc == nil {
		panic("")
	}
	return _authfunc
}

func SetCaller(f CallFunc) {
	_callfunc = f
}

func Caller() CallFunc {
	if _callfunc == nil {
		panic("")
	}
	return _callfunc
}
