package mapclient

import (
	"context"
	"fmt"
	"m3game/client"
	"m3game/proto/pb"
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
			dpb.File_map_proto.Services().Get(0),
			srcins,
			&pb.RouteSvc{
				EnvID:   srcins.EnvID,
				WorldID: srcins.WorldID,
				FuncID:  srcins.FuncID,
				IDStr:   util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.MapAppFuncID),
			},
		),
		opts: opts,
	}

	var err error
	if _client.conn, err = grpc.Dial(
		fmt.Sprintf("router://%s", util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.MapAppFuncID)),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(client.SendInteror()),
	); err != nil {
		return err
	} else {
		_client.MapSerClient = dpb.NewMapSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Meta
	dpb.MapSerClient
	conn *grpc.ClientConn
	opts []grpc.CallOption
}

func Move(ctx context.Context, name string, distance int32, opts ...grpc.CallOption) (string, int32, error) {
	var in dpb.Move_Req
	in.Name = name
	in.Distance = distance
	out, err := client.RPCCallRandom(_client, _client.Move, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", 0, err
	} else {
		return out.Name, out.Location, nil
	}
}
