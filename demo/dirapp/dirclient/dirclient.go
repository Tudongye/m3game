package dirclient

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
			dpb.File_dir_proto.Services().Get(0),
			srcins,
			&pb.RouteSvc{
				EnvID:   srcins.EnvID,
				WorldID: srcins.WorldID,
				FuncID:  srcins.FuncID,
				IDStr:   util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.DirAppFuncID),
			},
		),
		opts: opts,
	}
	var err error
	if _client.conn, err = grpc.Dial(
		fmt.Sprintf("router://%s", util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.DirAppFuncID)),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(client.SendInteror()),
	); err != nil {
		return err
	} else {
		_client.DirSerClient = dpb.NewDirSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Meta
	dpb.DirSerClient
	conn *grpc.ClientConn
	opts []grpc.CallOption
}

func Hello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in dpb.Hello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.Hello, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}
