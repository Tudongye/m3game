package uidcli

import (
	"context"
	"fmt"
	"m3game/plugins/log"
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
	if err := rpc.InjectionRPC(pb.File_uid_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC uid %s", err.Error())
	}
}

func New(srcapp meta.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	if env, world, _, _, err := srcapp.Parse(); err != nil {
		return nil, nil
	} else {
		dstsvc := meta.GenRouteSvc(env, world, proto.UidFuncID)
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
		_client.UidSerClient = pb.NewUidSerClient(_client.conn)
		return _client, nil
	}
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

type Client struct {
	client.Client
	pb.UidSerClient
	conn *grpc.ClientConn
}

func AllocRoleId(ctx context.Context, openid string, opts ...grpc.CallOption) (string, error) {
	var in pb.AllocRoleId_Req
	in.OpenId = openid
	out, err := client.RPCCallSingle(_client, _client.AllocRoleId, ctx, &in, opts...)
	if err != nil {
		return "", err
	} else {
		return out.RoleId, nil
	}
}

func AllocClubId(ctx context.Context, roleid string, opts ...grpc.CallOption) (string, error) {
	var in pb.AllocClubId_Req
	in.RoleId = roleid
	out, err := client.RPCCallSingle(_client, _client.AllocClubId, ctx, &in, opts...)
	if err != nil {
		return "", err
	} else {
		return out.ClubId, nil
	}
}
