package rolechclient

import (
	"context"
	"fmt"
	"m3game/proto/pb"
	"m3game/runtime/client"
	"m3game/util"

	dproto "m3game/demo/proto"
	dpb "m3game/demo/proto/pb"

	"google.golang.org/grpc"
)

var (
	_client *Client
)

func Init(srcins *pb.RouteIns, opts ...grpc.CallOption) error {
	_client = &Client{
		Meta: client.NewMeta(
			dpb.File_rolech_proto.Services().Get(0),
			srcins,
			&pb.RouteSvc{
				EnvID:   srcins.EnvID,
				WorldID: srcins.WorldID,
				FuncID:  srcins.FuncID,
				IDStr:   util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.RoleAppFuncID),
			},
		),
		opts: opts,
	}

	var err error
	if _client.conn, err = grpc.Dial(
		fmt.Sprintf("router://%s", util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.RoleAppFuncID)),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(client.SendInteror()),
	); err != nil {
		return err
	} else {
		_client.RoleChSerClient = dpb.NewRoleChSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Meta
	dpb.RoleChSerClient
	conn *grpc.ClientConn
	opts []grpc.CallOption
}

func TransChannel(ctx context.Context, msg *dpb.ChannelMsg, opts ...grpc.CallOption) error {
	var in dpb.TransChannel_Req
	in.Msg = msg
	_, err := client.RPCCallBroadCast(_client, _client.TransChannel, ctx, &in, append(opts, _client.opts...)...)
	return err
}
