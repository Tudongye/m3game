package client

import (
	"m3game/meta"

	"google.golang.org/grpc"
)

var (
	_clientInterceptors []grpc.UnaryClientInterceptor
)

func RegisterClientInterceptor(f grpc.UnaryClientInterceptor) {
	_clientInterceptors = append(_clientInterceptors, f)
}

func ClientInterceptors() []grpc.UnaryClientInterceptor {
	return _clientInterceptors
}

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
