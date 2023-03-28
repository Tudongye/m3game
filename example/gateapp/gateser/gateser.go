package gateser

import (
	"context"
	"m3game/example/proto/pb"
	"m3game/meta/metapb"
	"m3game/plugins/gate"
	"m3game/plugins/log"
	"m3game/runtime/rpc"
	"m3game/runtime/server/multi"

	"google.golang.org/grpc"
	gpb "google.golang.org/protobuf/proto"
)

func init() {
	if err := rpc.InjectionRPC(pb.File_gate_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC GateSer %s", err.Error())
	}
}

func New() *GateSer {
	return &GateSer{
		Server: multi.New("GateSer"),
	}
}

type GateSer struct {
	*multi.Server
	pb.UnimplementedGateSerServer
}

func (d *GateSer) SendToCli(ctx context.Context, in *pb.SendToCli_Req) (*pb.SendToCli_Rsp, error) {
	out := new(pb.SendToCli_Rsp)
	playerid := in.PlayerID
	if csconn := gate.GetConn(playerid); csconn != nil {
		csmsg := &metapb.CSMsg{}
		csmsg.Content, _ = gpb.Marshal(in)
		csconn.Send(ctx, csmsg)
	} else {
		log.Error("not find playerid %s", playerid)
	}
	return out, nil
}

func (s *GateSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterGateSerServer(t, s)
		return nil
	}
}
