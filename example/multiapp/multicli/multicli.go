package multicli

import (
	"context"
	"fmt"
	"m3game/runtime/client"
	"m3game/runtime/rpc"

	"m3game/example/proto"
	"m3game/example/proto/pb"

	"github.com/pkg/errors"

	"m3game/meta"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_multi_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("InjectionRPC Multi %s", err.Error()))
	}
}

func Init(srcapp meta.RouteApp, opts ...grpc.DialOption) error {
	if _client != nil {
		return nil
	}
	if env, world, _, _, err := srcapp.Parse(); err != nil {
		return nil
	} else {
		dstsvc := meta.GenRouteSvc(env, world, proto.MultiAppFuncID)
		_client = &Client{
			Client: client.New(srcapp, dstsvc),
		}
	}
	var err error
	target := fmt.Sprintf("router://%s", _client.DstSvc().String())
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`))
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(client.ClientInterceptors()...)))
	if _client.conn, err = grpc.Dial(target, opts...); err != nil {
		return errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.MultiSerClient = pb.NewMultiSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Client
	pb.MultiSerClient
	conn *grpc.ClientConn
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

func Hello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in pb.Hello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.Hello, ctx, &in, opts...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}

func TraceHello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in pb.TraceHello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.TraceHello, ctx, &in, opts...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}

func BreakHello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in pb.BreakHello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.BreakHello, ctx, &in, opts...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}
