package rolechserver

import (
	"context"
	"fmt"
	dpb "m3game/demo/proto/pb"
	"m3game/runtime/rpc"
	"m3game/runtime/server/async"

	"google.golang.org/grpc"
)

var (
	_channmsgs []*dpb.ChannelMsg
)

func AppendMsg(msg *dpb.ChannelMsg) {
	_channmsgs = append(_channmsgs, msg)
}

func GetMsg() []*dpb.ChannelMsg {
	return _channmsgs
}

func init() {
	if err := rpc.RegisterRPCSvc(dpb.File_rolech_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc RoleCh %s", err.Error()))
	}
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
