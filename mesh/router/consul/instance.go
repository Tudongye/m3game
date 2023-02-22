package consul

import (
	"m3game/mesh/router"
)

type Instance struct {
	host  string
	port  uint32
	idstr string
}

var (
	_ router.Instance = (*Instance)(nil)
)

func (i *Instance) GetHost() string {
	return i.host
}

func (i *Instance) GetPort() uint32 {
	return i.port
}

func (i *Instance) GetIDStr() string {
	return i.idstr
}
