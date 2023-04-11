package client

import "m3game/runtime/mesh"

type Client interface {
	SrcApp() mesh.RouteApp
	DstSvc() mesh.RouteSvc
}

type clientBase struct {
	srcapp mesh.RouteApp
	dstsvc mesh.RouteSvc
}

func New(srcapp mesh.RouteApp, dstsvc mesh.RouteSvc) Client {
	m := &clientBase{
		srcapp: srcapp,
		dstsvc: dstsvc,
	}
	return m
}

func (m *clientBase) SrcApp() mesh.RouteApp {
	return m.srcapp
}

func (m *clientBase) DstSvc() mesh.RouteSvc {
	return m.dstsvc
}
