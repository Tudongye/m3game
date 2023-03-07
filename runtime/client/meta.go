package client

import (
	"m3game/proto/pb"
)

type Meta interface {
	SrcIns() *pb.RouteIns
	DstSvc() *pb.RouteSvc
}

func NewMeta(srcins *pb.RouteIns, dstsvc *pb.RouteSvc) Meta {
	m := &defaultMeta{
		srcins: srcins,
		dstsvc: dstsvc,
	}
	return m
}

type defaultMeta struct {
	srcins *pb.RouteIns
	dstsvc *pb.RouteSvc
}

func (m *defaultMeta) SrcIns() *pb.RouteIns {
	return m.srcins
}
func (m *defaultMeta) DstSvc() *pb.RouteSvc {
	return m.dstsvc
}
