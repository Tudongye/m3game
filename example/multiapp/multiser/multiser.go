package multiser

import (
	"context"
	"fmt"
	"m3game/example/proto/pb"
	"m3game/plugins/log"
	"m3game/runtime/rpc"
	"m3game/runtime/server/multi"

	"google.golang.org/grpc"
)

func init() {
	if err := rpc.InjectionRPC(pb.File_multi_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC MultiSer %s", err.Error())
	}
}

func New() *MultiSer {
	return &MultiSer{
		Server: multi.New("MultiSer"),
	}
}

type MultiSer struct {
	*multi.Server
	pb.UnimplementedMultiSerServer
}

func (d *MultiSer) Hello(ctx context.Context, in *pb.Hello_Req) (*pb.Hello_Rsp, error) {
	out := new(pb.Hello_Rsp)
	out.Rsp = fmt.Sprintf("Hello , %s", in.Req)
	return out, nil
}

func (d *MultiSer) TraceHello(ctx context.Context, in *pb.TraceHello_Req) (*pb.TraceHello_Rsp, error) {
	out := new(pb.TraceHello_Rsp)
	out.Rsp = fmt.Sprintf("Hello , %s", in.Req)
	return out, nil
}

func (d *MultiSer) BreakHello(ctx context.Context, in *pb.BreakHello_Req) (*pb.BreakHello_Rsp, error) {
	out := new(pb.BreakHello_Rsp)
	out.Rsp = fmt.Sprintf("Hello , %s", in.Req)
	return out, nil
}

func (s *MultiSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterMultiSerServer(t, s)
		return nil
	}
}
