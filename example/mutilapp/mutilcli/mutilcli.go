package mutilcli

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
	if err := rpc.RegisterRPCSvc(pb.File_mutil_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc Mutil %s", err.Error()))
	}
}

func Init(srcins *mpb.RouteIns, opts ...grpc.CallOption) error {
	if _client != nil {
		return nil
	}
	dstsvc := util.RouteIns2Svc(srcins, proto.MutilAppFuncID)
	_client = &Client{
		Meta: client.NewMeta(
			srcins,
			dstsvc,
		),
		opts: opts,
	}
	var err error
	target := fmt.Sprintf("router://%s", dstsvc.IDStr)
	if _client.conn, err = grpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(transport.ClientInterceptors()...)),
	); err != nil {
		return errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.MutilSerClient = pb.NewMutilSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Meta
	pb.MutilSerClient
	conn *grpc.ClientConn
	opts []grpc.CallOption
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

func Hello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in pb.Hello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.Hello, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}

func TraceHello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in pb.TraceHello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.TraceHello, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}

func BreakHello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in pb.BreakHello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.BreakHello, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}
