package consul

import (
	"m3game/plugins/router"
)

type Instance struct {
	host  string
	port  uint32
	idstr string
}

func newInstance(host string, port int, idstr string) *Instance {
	return &Instance{
		host:  host,
		port:  uint32(port),
		idstr: idstr,
	}
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
