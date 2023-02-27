package roleclient

import (
	"context"
	"fmt"
	"m3game/client"
	"m3game/proto"
	"m3game/proto/pb"
	"m3game/runtime/transport"
	"m3game/util"

	dpb "m3game/demo/proto/pb"

	dproto "m3game/demo/proto"

	"google.golang.org/grpc"
)

var (
	_instance     *Client
	_rolemetaskey = "_rolemetaskey"
)

type Opt func(*Client)

func RoleClient() *Client {
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
		_instance.client = dpb.NewRoleSerClient(_instance.conn)
		return nil
	}
}

func SendInterFunc(s *transport.Sender) error {
	s.Metas().Set(proto.META_CLIENT, _instance.Client)
	if m := s.Ctx().Value(_rolemetaskey); m != nil {
		metas := m.(map[string]string)
		for k, v := range metas {
			s.Metas().Set(k, v)
		}
	}
	return client.SendInterFunc(s)
}

type Client struct {
	conn   *grpc.ClientConn
	srcins *pb.RouteIns
	dstsvc *pb.RouteSvc
	client dpb.RoleSerClient

	Client string
}

func (c *Client) Register(ctx context.Context, roleid string, name string) (string, error) {
	var in dpb.Register_Req
	in.RouteHead = client.CreateRouteHead_Hash(c.srcins, c.dstsvc, roleid)
	metas := make(map[string]string)
	metas[proto.META_ACTORID] = roleid
	metas[proto.META_CREATE_ACTORID] = proto.META_FLAG_TRUE
	ctx = context.WithValue(ctx, _rolemetaskey, metas)
	in.Name = name
	if out, err := c.client.Register(ctx, &in); err != nil {
		return "", err
	} else {
		return out.RoleID, nil
	}
}

func (c *Client) Login(ctx context.Context, roleid string) (string, string, error) {
	var in dpb.Login_Req
	in.RouteHead = client.CreateRouteHead_Hash(c.srcins, c.dstsvc, roleid)
	metas := make(map[string]string)
	metas[proto.META_ACTORID] = roleid
	metas[proto.META_CREATE_ACTORID] = proto.META_FLAG_TRUE
	ctx = context.WithValue(ctx, _rolemetaskey, metas)
	if out, err := c.client.Login(ctx, &in); err != nil {
		return "", "", err
	} else {
		return out.Name, out.Tips, nil
	}
}

func (c *Client) ModifyName(ctx context.Context, roleid string, name string) (string, error) {
	var in dpb.ModifyName_Req
	in.RouteHead = client.CreateRouteHead_Hash(c.srcins, c.dstsvc, roleid)
	in.NewName = name
	metas := make(map[string]string)
	metas[proto.META_ACTORID] = roleid
	ctx = context.WithValue(ctx, _rolemetaskey, metas)
	if out, err := c.client.ModifyName(ctx, &in); err != nil {
		return "", err
	} else {
		return out.Name, nil
	}
}

func (c *Client) GetName(ctx context.Context, roleid string) (string, error) {
	var in dpb.GetName_Req
	in.RouteHead = client.CreateRouteHead_Hash(c.srcins, c.dstsvc, roleid)
	metas := make(map[string]string)
	metas[proto.META_ACTORID] = roleid
	ctx = context.WithValue(ctx, _rolemetaskey, metas)
	if out, err := c.client.GetName(ctx, &in); err != nil {
		return "", err
	} else {
		return out.Name, nil
	}
}

func (c *Client) MoveRole(ctx context.Context, roleid string, distance int32) (int32, string, error) {
	var in dpb.MoveRole_Req
	in.RouteHead = client.CreateRouteHead_Hash(c.srcins, c.dstsvc, roleid)
	in.Distance = distance
	metas := make(map[string]string)
	metas[proto.META_ACTORID] = roleid
	ctx = context.WithValue(ctx, _rolemetaskey, metas)
	if out, err := c.client.MoveRole(ctx, &in); err != nil {
		return 0, "", err
	} else {
		return out.Location, out.LocateName, nil
	}
}

func (c *Client) PostChannel(ctx context.Context, roleid string, content string) error {
	var in dpb.PostChannel_Req
	in.RouteHead = client.CreateRouteHead_Hash(c.srcins, c.dstsvc, roleid)
	in.Content = content
	metas := make(map[string]string)
	metas[proto.META_ACTORID] = roleid
	ctx = context.WithValue(ctx, _rolemetaskey, metas)
	if _, err := c.client.PostChannel(ctx, &in); err != nil {
		return err
	} else {
		return nil
	}
}

func (c *Client) PullChannel(ctx context.Context, roleid string) ([]*dpb.ChannelMsg, error) {
	var in dpb.PullChannel_Req
	in.RouteHead = client.CreateRouteHead_Hash(c.srcins, c.dstsvc, roleid)
	metas := make(map[string]string)
	metas[proto.META_ACTORID] = roleid
	ctx = context.WithValue(ctx, _rolemetaskey, metas)
	if out, err := c.client.PullChannel(ctx, &in); err != nil {
		return nil, err
	} else {
		return out.Msgs, nil
	}
}
