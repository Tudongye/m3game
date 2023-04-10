package client

import (
	"m3game/meta"
)

type Client interface {
	SrcApp() meta.RouteApp
	DstSvc() meta.RouteSvc
}

type clientBase struct {
	srcapp meta.RouteApp
	dstsvc meta.RouteSvc
}

func New(srcapp meta.RouteApp, dstsvc meta.RouteSvc) Client {
	m := &clientBase{
		srcapp: srcapp,
		dstsvc: dstsvc,
	}
	return m
}

func (m *clientBase) SrcApp() meta.RouteApp {
	return m.srcapp
}

func (m *clientBase) DstSvc() meta.RouteSvc {
	return m.dstsvc
}
