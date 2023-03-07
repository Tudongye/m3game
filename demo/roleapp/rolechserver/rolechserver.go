package rolechserver

import (
	"context"
	dpb "m3game/demo/proto/pb"
	"m3game/runtime/server/async"

	"google.golang.org/grpc"
)

var (
	_map       map[string]int32
	_channmsgs []*dpb.ChannelMsg
)

func AppendMsg(msg *dpb.ChannelMsg) {
	_channmsgs = append(_channmsgs, msg)
}

func GetMsg() []*dpb.ChannelMsg {
	return _channmsgs
}

func init() {
	_map = make(map[string]int32)
}

func New() *RoleChSer {
	return &RoleChSer{
		Server: async.New("RoleChSer"),
	}
}

type RoleChSer struct {
	*async.Server
	dpb.UnimplementedRoleChSerServer
}

func (d *RoleChSer) TransChannel(ctx context.Context, in *dpb.TransChannel_Req) (*dpb.TransChannel_Rsp, error) {
	out := new(dpb.TransChannel_Rsp)
	AppendMsg(in.Msg)
	return out, nil
}

func (s *RoleChSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		dpb.RegisterRoleChSerServer(t, s)
		return nil
	}
}
