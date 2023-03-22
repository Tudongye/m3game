package simpleser

import (
	"context"
	"fmt"
	"m3game/example/proto/pb"
	"m3game/runtime/rpc"
	"m3game/runtime/server/multi"

	"google.golang.org/grpc"
)

func init() {
	// 注册RPC信息到框架层
	if err := rpc.InjectionRPC(pb.File_simple_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("InjectionRPC SimpleSer %s", err.Error()))
	}
}

func New() *SimpleSer {
	return &SimpleSer{
		Server: multi.New("SimpleSer"), // 以MultiSer为基础构建SimpleSer
	}
}

type SimpleSer struct {
	*multi.Server
	pb.UnimplementedSimpleSerServer
}

// 实现HelloWorld接口
func (d *SimpleSer) HelloWorld(ctx context.Context, in *pb.HelloWorld_Req) (*pb.HelloWorld_Rsp, error) {
	out := new(pb.HelloWorld_Rsp)
	out.Rsp = fmt.Sprintf("HelloWorld , %s", in.Req)
	return out, nil
}

// 将SimpleSer注册到grpcser
func (s *SimpleSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterSimpleSerServer(t, s)
		return nil
	}
}
