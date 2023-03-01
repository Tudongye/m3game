package proto

import (
	"m3game/proto/pb"
)

type M3Pkg interface {
	GetRouteHead() *pb.RouteHead
}

type M3Metas struct {
	Metas map[string]string
}

func NewM3Metas() *M3Metas {
	return &M3Metas{
		Metas: make(map[string]string),
	}
}

func (m *M3Metas) Decode(pm *pb.Metas) {
	if pm == nil {
		return
	}
	for _, meta := range pm.Metas {
		m.Metas[meta.Key] = meta.Value
	}
}

func (m *M3Metas) Encode() *pb.Metas {
	var pm pb.Metas
	for k, v := range m.Metas {
		pm.Metas = append(pm.Metas, &pb.Meta{
			Key:   k,
			Value: v,
		})
	}
	return &pm
}

func (m *M3Metas) Set(key string, value string) {
	m.Metas[key] = value
}

func (m *M3Metas) Get(key string) (string, bool) {
	v, ok := m.Metas[key]
	return v, ok
}
