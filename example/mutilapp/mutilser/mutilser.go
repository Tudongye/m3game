package mutilser

import (
	"context"
	"fmt"
	"m3game/example/proto"
	"m3game/example/proto/pb"
	"m3game/runtime/rpc"
	"m3game/runtime/server"
	"m3game/runtime/server/mutil"

	"google.golang.org/grpc"
)

func init() {
	if err := rpc.RegisterRPCSvc(pb.File_mutil_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc MutilSer %s", err.Error()))
	}
}

func New() *MutilSer {
	return &MutilSer{
		Server: mutil.New("MutilSer"),
	}
}

type MutilSer struct {
	*mutil.Server
	pb.UnimplementedMutilSerServer
}

func (d *MutilSer) Hello(ctx context.Context, in *pb.Hello_Req) (*pb.Hello_Rsp, error) {
	out := new(pb.Hello_Rsp)
	sctx := server.ParseContext(ctx)
	out.Rsp = fmt.Sprintf("Hello , %s", in.Req)
	if sctx != nil {
		if _, ok := sctx.Reciver().Metas().Get(proto.RHMeta_Client); ok {
			out.Rsp = fmt.Sprintf("Hello Client , %s", in.Req)
		}
	}
	return out, nil
}

func (d *MutilSer) TraceHello(ctx context.Context, in *pb.TraceHello_Req) (*pb.TraceHello_Rsp, error) {
	out := new(pb.TraceHello_Rsp)
	sctx := server.ParseContext(ctx)
	out.Rsp = fmt.Sprintf("TraceHello , %s", in.Req)
	if sctx != nil {
		if _, ok := sctx.Reciver().Metas().Get(proto.RHMeta_Client); ok {
			out.Rsp = fmt.Sprintf("TraceHello Client , %s", in.Req)
		}
	}
	return out, nil
}

func (d *MutilSer) BreakHello(ctx context.Context, in *pb.BreakHello_Req) (*pb.BreakHello_Rsp, error) {
	out := new(pb.BreakHello_Rsp)
	sctx := server.ParseContext(ctx)
	out.Rsp = fmt.Sprintf("BreakHello , %s", in.Req)
	if sctx != nil {
		if _, ok := sctx.Reciver().Metas().Get(proto.RHMeta_Client); ok {
			out.Rsp = fmt.Sprintf("BreakHello Client , %s", in.Req)
		}
	}
	return out, nil
}

func (s *MutilSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterMutilSerServer(t, s)
		return nil
	}
}
