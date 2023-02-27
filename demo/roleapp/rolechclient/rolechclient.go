package rolechclient

import (
	"context"
	"fmt"
	"m3game/client"
	"m3game/proto"
	"m3game/proto/pb"
	"m3game/runtime/transport"
	"m3game/util"

	dproto "m3game/demo/proto"
	dpb "m3game/demo/proto/pb"

	"google.golang.org/grpc"
)

var (
	_instance *Client
)

type Opt func(*Client)

func RoleChClient() *Client {
	return _instance
}

func Init(srcins *pb.RouteIns, opts ...Opt) error {
	_instance = &Client{
		srcins: srcins,
		dstsvc: &pb.RouteSvc{
			EnvID:   srcins.EnvID,
			WorldID: srcins.WorldID,
			FuncID:  srcins.FuncID,
			IDStr:   util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.RoleAppFuncID),
		},
	}
	for _, opt := range opts {
		opt(_instance)
	}
	var err error
	if _instance.conn, err = grpc.Dial(
		fmt.Sprintf("router://%s", util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.RoleAppFuncID)),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(transport.SendInteror(SendInterFunc)),
	); err != nil {
		return err
	} else {
		_instance.client = dpb.NewRoleChSerClient(_instance.conn)
		return nil
	}
}

func SendInterFunc(s *transport.Sender) error {
	s.Metas().Set(proto.META_CLIENT, _instance.Client)
	return client.SendInterFunc(s)
}

type Client struct {
	conn   *grpc.ClientConn
	srcins *pb.RouteIns
	dstsvc *pb.RouteSvc
	client dpb.RoleChSerClient

	Client string
}

func (c *Client) TransChannel(ctx context.Context, msg *dpb.ChannelMsg) error {
	var in dpb.TransChannel_Req
	in.RouteHead = client.CreateRouteHead_Broad(c.srcins, c.dstsvc)
	in.Msg = msg
	if _, err := c.client.TransChannel(ctx, &in); err != nil {
		return err
	} else {
		return nil
	}
}
