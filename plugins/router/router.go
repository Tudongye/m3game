package router

import "m3game/runtime/plugin"

var (
	_router Router
)

type Instance interface {
	GetHost() string
	GetPort() uint32
	GetIDStr() string
}

type Router interface {
	plugin.PluginIns
	Register(app string, svc string, addr string) error
	Deregister(app string, svc string) error
	GetAllInstances(svcid string) ([]Instance, error)
}

func Set(r Router) {
	if _router != nil {
		panic("Router Only One")
	}
	_router = r
}

func Get() Router {
	if _router == nil {
		panic("Router Mush Have One")
	}
	return _router
}

func Register(app string, svc string, addr string) error {
	return Get().Register(app, svc, addr)
}

func Deregister(app string, svc string) error {
	return Get().Deregister(app, svc)
}

func GetAllInstances(svcid string) ([]Instance, error) {
	return Get().GetAllInstances(svcid)
}
