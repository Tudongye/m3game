package dirclient

import (
	"context"
	"fmt"
	"m3game/client"
	"m3game/proto"
	"m3game/proto/pb"
	"m3game/runtime/transport"
	"m3game/util"

	"m3game/demo/dirapp"
	dpb "m3game/demo/proto/pb"

	"google.golang.org/grpc"
)

var (
	_instance *Client
)

type Opt func(*Client)

func DirClient() *Client {
	return _instance
}

func Init(srcins *pb.RouteIns, opts ...Opt) error {
	_instance = &Client{
		srcins: srcins,
		dstsvc: &pb.RouteSvc{
			EnvID:   srcins.EnvID,
			WorldID: srcins.WorldID,
			FuncID:  srcins.FuncID,
			IDStr:   util.SvcID2Str(srcins.EnvID, srcins.WorldID, dirapp.AppFuncID),
		},
	}
	for _, opt := range opts {
		opt(_instance)
	}
	var err error
	if _instance.conn, err = grpc.Dial(
		fmt.Sprintf("router://%s", util.SvcID2Str(srcins.EnvID, srcins.WorldID, dirapp.AppFuncID)),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(transport.SendInteror(SendInterFunc)),
	); err != nil {
		return err
	} else {
		_instance.client = dpb.NewDirSerClient(_instance.conn)
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
	client dpb.DirSerClient

	Client string
}

func (c *Client) Hello(ctx context.Context, hellostr string) (string, error) {
	var in dpb.Hello_Req
	in.RouteHead = client.CreateRouteHead_Random(c.srcins, c.dstsvc)
	in.Req = hellostr
	if out, err := c.client.Hello(ctx, &in); err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}
