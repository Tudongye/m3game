package simplecli

import (
	"context"
	"fmt"
	mpb "m3game/proto/pb"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"m3game/runtime/transport"
	"m3game/util"

	"m3game/example/proto"
	"m3game/example/proto/pb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"

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

func Init(srcins *mpb.RouteIns, opts ...grpc.CallOption) error {
	if _client != nil {
		return nil
	}
	// SimpleApp的服务名
	dstsvc := util.RouteIns2Svc(srcins, proto.SimpleAppFuncID)
	_client = &Client{
		Meta: client.NewMeta(
			srcins,
			dstsvc,
		),
		opts: opts,
	}
	var err error
	// 使用M3 自定义的Router负载均衡器（寻址）
	target := fmt.Sprintf("router://%s", dstsvc.IDStr)
	if _client.conn, err = grpc.Dial(
		target,
		grpc.WithInsecure(),
		// 注入框架层拦截器
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(transport.ClientInterceptors()...)),
	); err != nil {
		return errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.SimpleSerClient = pb.NewSimpleSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Meta
	pb.SimpleSerClient
	conn *grpc.ClientConn
	opts []grpc.CallOption
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

func HelloWorld(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in pb.HelloWorld_Req
	in.Req = hellostr
	// 填充路由信息，并调用RPC接口
	out, err := client.RPCCallRandom(_client, _client.HelloWorld, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}
