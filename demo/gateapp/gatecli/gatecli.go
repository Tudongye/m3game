package gatecli

import (
	"context"
	"fmt"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"time"

	"m3game/demo/proto"
	"m3game/demo/proto/pb"

	"github.com/pkg/errors"

	"m3game/meta"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_gate_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("InjectionRPC Gate %s", err.Error()))
	}
}

func New(srcapp meta.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	if env, world, _, _, err := srcapp.Parse(); err != nil {
		return nil, nil
	} else {
		dstsvc := meta.GenRouteSvc(env, world, proto.GateFuncID)
		_client = &Client{
			Client: client.New(srcapp, dstsvc),
		}
	}
	var err error
	target := fmt.Sprintf("router://%s", _client.DstSvc().String())
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(client.ClientInterceptors()...)),
		grpc.WithTimeout(time.Second*10),
	)
	if _client.conn, err = grpc.Dial(target, opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.GateSerClient = pb.NewGateSerClient(_client.conn)
		return _client, nil
	}
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

type Client struct {
	client.Client
	pb.GateSerClient
	conn *grpc.ClientConn
}

func SendToCli(ctx context.Context, roleid string, ntymsg *pb.NtyMsg, dstapp meta.RouteApp, opts ...grpc.CallOption) error {
	var in pb.SendToCli_Req
	in.RoleId = roleid
	in.NtyMsg = ntymsg
	_, err := client.RPCCallP2P(_client, _client.SendToCli, ctx, &in, dstapp, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}
