package asyncser

import (
	"context"
	"m3game/example/proto/pb"
	"m3game/plugins/log"
	"m3game/runtime/rpc"
	"m3game/runtime/server/async"
	"time"

	"google.golang.org/grpc"
)

func init() {
	if err := rpc.InjectionRPC(pb.File_async_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC AsyncSer %s", err.Error())
	}
}

func New() *AsyncSer {
	return &AsyncSer{
		Server: async.New("AsyncSer"),
	}
}

type AsyncSer struct {
	*async.Server
	pb.UnimplementedAsyncSerServer
}

func (d *AsyncSer) TransChannel(ctx context.Context, in *pb.TransChannel_Req) (*pb.TransChannel_Rsp, error) {
	out := new(pb.TransChannel_Rsp)
	AppendMsg(in.Msg)
	return out, nil
}

func (d *AsyncSer) SSPullChannel(ctx context.Context, in *pb.SSPullChannel_Req) (*pb.SSPullChannel_Rsp, error) {
	out := new(pb.SSPullChannel_Rsp)
	out.Msgs = GetMsg()
	time.Sleep(5 * time.Second)
	return out, nil
}

func (s *AsyncSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterAsyncSerServer(t, s)
		return nil
	}
}
