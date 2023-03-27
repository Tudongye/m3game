package router

import (
	"fmt"
	"m3game/plugins/log"
	"m3game/runtime/plugin"
)

var (
	_router Router
)

type Ins interface {
	GetHost() string
	GetPort() uint32
	GetIDStr() string
	GetMeta(string) (string, bool)
}

type Router interface {
	plugin.PluginIns
	Register(app string, svc string, addr string, meta map[string]string) error
	Deregister(app string, svc string) error
	GetAllInstances(svcid string) ([]Ins, error)
}

func New(me Router) (Router, error) {
	if _router != nil {
		log.Fatal("Metric Only One")
		return nil, fmt.Errorf("Metric is newed %s", me.Factory().Name())
	}
	_router = me
	return _router, nil
}

func Instance() Router {
	if _router == nil {
		log.Fatal("Router not newd")
		return nil
	}
	return _router
}

func Register(app string, svc string, addr string, meta map[string]string) error {
	return Instance().Register(app, svc, addr, meta)
}

func Deregister(app string, svc string) error {
	return Instance().Deregister(app, svc)
}

func GetAllInstances(svcid string) ([]Ins, error) {
	return Instance().GetAllInstances(svcid)
}
