package uidcli

import (
	"context"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/rpc"

	"m3game/demo/proto"
	"m3game/demo/proto/pb"

	"github.com/pkg/errors"

	"m3game/meta"

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
	dstsvc := meta.GenDstRouteSvc(srcapp, proto.UidFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.UidSerClient = pb.NewUidSerClient(_client.conn)
		return _client, nil
	}
}

func Conn() grpc.ClientConnInterface {
	return _client.conn
}

type Client struct {
	client.Client
	pb.UidSerClient
	conn grpc.ClientConnInterface
}

func AllocRoleId(ctx context.Context, openid string, opts ...grpc.CallOption) (int64, error) {
	var in pb.AllocRoleId_Req
	in.OpenId = openid
	out, err := client.RPCCallSingle(_client, _client.AllocRoleId, ctx, &in, opts...)
	if err != nil {
		return 0, err
	} else {
		return out.RoleId, nil
	}
}

func AllocClubId(ctx context.Context, roleid int64, opts ...grpc.CallOption) (int64, error) {
	var in pb.AllocClubId_Req
	in.RoleId = roleid
	out, err := client.RPCCallSingle(_client, _client.AllocClubId, ctx, &in, opts...)
	if err != nil {
		return 0, err
	} else {
		return out.ClubId, nil
	}
}
