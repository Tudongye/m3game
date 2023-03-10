package client

import (
	"m3game/proto"
	"m3game/runtime/transport"

	"google.golang.org/grpc"
)

type M3CallOption interface {
	grpc.CallOption
	Filter() func(*transport.Sender) error
}

type M3Option struct {
	grpc.EmptyCallOption
	F func(*transport.Sender) error
}

func (m *M3Option) Filter() func(*transport.Sender) error {
	return m.F
}

var _ M3CallOption = (*M3Option)(nil)
var _ grpc.CallOption = (*M3Option)(nil)

// 请求是否来自客户端
func GenMetaClientOption(v string) *M3Option {
	return &M3Option{
		F: func(s *transport.Sender) error {
			s.Metas().Set(proto.META_CLIENT, v)
			return nil
		},
	}
}

// Actor Server :  actorid
func GenMetaActorIDOption(v string) *M3Option {
	return &M3Option{
		F: func(s *transport.Sender) error {
			s.Metas().Set(proto.META_ACTORID, v)
			return nil
		},
	}
}

// Actor Server :  bool create actor
func GenMetaCreateActorOption(v string) *M3Option {
	return &M3Option{
		F: func(s *transport.Sender) error {
			s.Metas().Set(proto.META_CREATE_ACTOR, v)
			return nil
		},
	}
}

// Actor Server :  bool create actor
func GenMetaOption(k string, v string) *M3Option {
	return &M3Option{
		F: func(s *transport.Sender) error {
			s.Metas().Set(k, v)
			return nil
		},
	}
}
