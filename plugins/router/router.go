package router

import (
	"m3game/meta/errs"
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
	Register(app string, svc string, host string, port int, meta map[string]string, livef func(app string, svc string) bool) error
	Deregister(app string, svc string) error
	GetAllInstances(svcid string) ([]Ins, error)
}

func New(me Router) (Router, error) {
	if _router != nil {
		log.Fatal("Metric Only One")
		return nil, errs.RouterInsHasNewed.New("Metric is newed %s", me.Factory().Name())
	}
	_router = me
	return _router, nil
}

func Instance() Router {
	if _router == nil {
		log.Error("Router not newd")
		return nil
	}
	return _router
}
