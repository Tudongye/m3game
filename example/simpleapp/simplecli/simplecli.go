package simplecli

import (
	"context"
	"fmt"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"m3game/runtime/transport"

	"m3game/example/proto"
	"m3game/example/proto/pb"

	"m3game/meta"

	"github.com/pkg/errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	// 注册RPC信息到框架层
	if err := rpc.RegisterRPCSvc(pb.File_simple_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc SimpleSer %s", err.Error()))
	}
}

func Init(srcapp meta.RouteApp, opts ...grpc.DialOption) error {
	if _client != nil {
		return nil
	}
	if env, world, _, _, err := srcapp.Parse(); err != nil {
		return nil
	} else {
		dstsvc := meta.GenRouteSvc(env, world, proto.SimpleAppFuncID)
		_client = &Client{
			Client: client.New(srcapp, dstsvc),
		}
	}
	var err error
	target := fmt.Sprintf("router://%s", _client.DstSvc().String())
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`))
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(transport.ClientInterceptors()...)))
	if _client.conn, err = grpc.Dial(target, opts...); err != nil {
		return errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.SimpleSerClient = pb.NewSimpleSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Client
	pb.SimpleSerClient
	conn *grpc.ClientConn
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

func HelloWorld(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in pb.HelloWorld_Req
	in.Req = hellostr
	// 填充路由信息，并调用RPC接口
	out, err := client.RPCCallRandom(_client, _client.HelloWorld, ctx, &in, opts...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}
