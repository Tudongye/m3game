package roleclient

import (
	"context"
	"fmt"
	"m3game/proto"
	"m3game/proto/pb"
	"m3game/runtime/client"
	"m3game/util"

	dpb "m3game/demo/proto/pb"

	dproto "m3game/demo/proto"

	"google.golang.org/grpc"
)

var (
	_client *Client
)

func Init(srcins *pb.RouteIns, opts ...grpc.CallOption) error {
	_client = &Client{
		Meta: client.NewMeta(
			dpb.File_role_proto.Services().Get(0),
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
		_client.RoleSerClient = dpb.NewRoleSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Meta
	dpb.RoleSerClient
	conn *grpc.ClientConn
	opts []grpc.CallOption
}

func Register(ctx context.Context, roleid string, name string, opts ...grpc.CallOption) (string, error) {
	var in dpb.Register_Req
	in.Name = name
	opts = append(opts, client.GenMetaActorIDOption(roleid))
	opts = append(opts, client.GenMetaCreateActorOption(proto.META_FLAG_TRUE))
	out, err := client.RPCCallRandom(_client, _client.Register, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.RoleID, nil
	}
}

func Login(ctx context.Context, roleid string, opts ...grpc.CallOption) (string, string, error) {
	var in dpb.Login_Req
	in.RoleID = roleid
	opts = append(opts, client.GenMetaActorIDOption(roleid))
	opts = append(opts, client.GenMetaCreateActorOption(proto.META_FLAG_TRUE))
	out, err := client.RPCCallHash(_client, _client.Login, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", "", err
	} else {
		return out.Name, out.Tips, nil
	}
}

func ModifyName(ctx context.Context, roleid string, name string, opts ...grpc.CallOption) (string, error) {
	var in dpb.ModifyName_Req
	in.RoleID = roleid
	in.NewName = name
	opts = append(opts, client.GenMetaActorIDOption(roleid))
	out, err := client.RPCCallHash(_client, _client.ModifyName, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Name, nil
	}
}

func GetName(ctx context.Context, roleid string, opts ...grpc.CallOption) (string, error) {
	var in dpb.GetName_Req
	in.RoleID = roleid
	opts = append(opts, client.GenMetaActorIDOption(roleid))
	out, err := client.RPCCallHash(_client, _client.GetName, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Name, nil
	}
}

func MoveRole(ctx context.Context, roleid string, distance int32, opts ...grpc.CallOption) (int32, string, error) {
	var in dpb.MoveRole_Req
	in.RoleID = roleid
	in.Distance = distance
	opts = append(opts, client.GenMetaActorIDOption(roleid))
	out, err := client.RPCCallHash(_client, _client.MoveRole, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return 0, "", err
	} else {
		return out.Location, out.LocateName, nil
	}
}

func PostChannel(ctx context.Context, roleid string, content string, opts ...grpc.CallOption) error {
	var in dpb.PostChannel_Req
	in.RoleID = roleid
	in.Content = content
	opts = append(opts, client.GenMetaActorIDOption(roleid))
	_, err := client.RPCCallHash(_client, _client.PostChannel, ctx, &in, append(opts, _client.opts...)...)
	return err
}

func PullChannel(ctx context.Context, roleid string, opts ...grpc.CallOption) ([]*dpb.ChannelMsg, error) {
	var in dpb.PullChannel_Req
	in.RoleID = roleid
	opts = append(opts, client.GenMetaActorIDOption(roleid))
	out, err := client.RPCCallHash(_client, _client.PullChannel, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return nil, err
	} else {
		return out.Msgs, nil
	}
}
