package simplecli

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
	// 注册RPC信息到框架层
	if err := rpc.InjectionRPC(pb.File_simple_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC SimpleSer %s", err.Error())
	}
}

func New(srcapp mesh.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := mesh.GenDstRouteSvc(srcapp, proto.SimpleAppFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.SimpleSerClient = pb.NewSimpleSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.SimpleSerClient
	conn grpc.ClientConnInterface
}

func Conn() grpc.ClientConnInterface {
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
