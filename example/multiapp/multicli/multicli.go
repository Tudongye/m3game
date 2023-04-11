package multicli

import (
	"context"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/mesh"
	"m3game/runtime/rpc"

	"m3game/example/proto"
	"m3game/example/proto/pb"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_multi_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC Multi %s", err.Error())
	}
}

func New(srcapp mesh.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := mesh.GenDstRouteSvc(srcapp, proto.MultiAppFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.MultiSerClient = pb.NewMultiSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.MultiSerClient
	conn grpc.ClientConnInterface
}

func Conn() grpc.ClientConnInterface {
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
