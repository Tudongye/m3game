package mapclient

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

func MapClient() *Client {
	return _instance
}

func Init(srcins *pb.RouteIns, opts ...Opt) error {
	_instance = &Client{
		srcins: srcins,
		dstsvc: &pb.RouteSvc{
			EnvID:   srcins.EnvID,
			WorldID: srcins.WorldID,
			FuncID:  srcins.FuncID,
			IDStr:   util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.MapAppFuncID),
		},
	}
	for _, opt := range opts {
		opt(_instance)
	}
	var err error
	if _instance.conn, err = grpc.Dial(
		fmt.Sprintf("router://%s", util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.MapAppFuncID)),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(transport.SendInteror(SendInterFunc)),
	); err != nil {
		return err
	} else {
		_instance.client = dpb.NewMapSerClient(_instance.conn)
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
	client dpb.MapSerClient

	Client string
}

func (c *Client) Move(ctx context.Context, name string, distance int32) (string, int32, error) {
	var in dpb.Move_Req
	in.RouteHead = client.CreateRouteHead_Random(c.srcins, c.dstsvc)
	in.Name = name
	in.Distance = distance
	if out, err := c.client.Move(ctx, &in); err != nil {
		return "", 0, err
	} else {
		return out.Name, out.Location, nil
	}
}
