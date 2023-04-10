package onlinecli

import (
	"context"
	"fmt"
	"m3game/plugins/log"
	"m3game/plugins/transport"
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
	if err := rpc.InjectionRPC(pb.File_online_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC online %s", err.Error())
	}
}

func New(srcapp meta.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := meta.GenDstRouteSvc(srcapp, proto.OnlineFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}
	var err error
	target := fmt.Sprintf("router://%s", _client.DstSvc().String())
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(transport.Instance().ClientInterceptors()...)),
		grpc.WithTimeout(time.Second*10),
	)
	if _client.conn, err = grpc.Dial(target, opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.OnlineSerClient = pb.NewOnlineSerClient(_client.conn)
		return _client, nil
	}
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

type Client struct {
	client.Client
	pb.OnlineSerClient
	conn *grpc.ClientConn
}

func OnlineCreate(ctx context.Context, roleid int64, appid string, opts ...grpc.CallOption) error {
	var in pb.OnlineCreate_Req
	in.RoleId = roleid
	in.AppId = appid
	_, err := client.RPCCallSingle(_client, _client.OnlineCreate, ctx, &in, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func OnlineDelete(ctx context.Context, roleid int64, appid string, opts ...grpc.CallOption) error {
	var in pb.OnlineDelete_Req
	in.RoleId = roleid
	in.AppId = appid
	_, err := client.RPCCallSingle(_client, _client.OnlineDelete, ctx, &in, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func OnlineRead(ctx context.Context, roleid int64, opts ...grpc.CallOption) (string, error) {
	var in pb.OnlineRead_Req
	in.RoleId = roleid
	out, err := client.RPCCallSingle(_client, _client.OnlineRead, ctx, &in, opts...)
	if err != nil {
		return "", err
	} else {
		return out.AppId, nil
	}
}
