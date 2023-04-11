package consul

import (
	"m3game/plugins/router"
)

type Instance struct {
	host  string
	port  uint32
	idstr string
	meta  map[string]string
}

func newInstance(host string, port int, idstr string, meta map[string]string) *Instance {
	return &Instance{
		host:  host,
		port:  uint32(port),
		idstr: idstr,
		meta:  meta,
	}
}

var (
	_ router.Ins = (*Instance)(nil)
)

func (i *Instance) GetHost() string {
	return i.host
}

func (i *Instance) GetPort() uint32 {
	return i.port
}

func (i *Instance) GetAppID() string {
	return i.idstr
}

func (i *Instance) GetMeta(k string) (string, bool) {
	v, ok := i.meta[k]
	return v, ok
}
