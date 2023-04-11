package clubcli

import (
	"context"
	"fmt"
	"m3game/demo/clubapp/club"
	"m3game/demo/proto"
	"m3game/demo/proto/pb"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/mesh"
	"m3game/runtime/rpc"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_club_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC Club %s", err.Error())
	}
}

func New(srcapp mesh.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := mesh.GenDstRouteSvc(srcapp, proto.ClubFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return _client, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.ClubSerClient = pb.NewClubSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.ClubSerClient
	conn grpc.ClientConnInterface
}

func Conn() grpc.ClientConnInterface {
	return _client.conn
}

func ClubCreate(ctx context.Context, clubid int64, roleid int64, opts ...grpc.CallOption) error {
	var in pb.ClubCreate_Req
	in.SlotId = club.Club2SlotId(clubid)
	in.ClubId = clubid
	in.RoleId = roleid
	md := metadata.Pairs(club.ClubIdMetaKey, fmt.Sprintf("%d", in.SlotId))
	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err := client.RPCCallHash(_client, _client.ClubCreate, ctx, &in, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ClubDelete(ctx context.Context, clubid int64, roleid int64, opts ...grpc.CallOption) error {
	var in pb.ClubDelete_Req
	in.SlotId = club.Club2SlotId(clubid)
	in.ClubId = clubid
	in.RoleId = roleid
	md := metadata.Pairs(club.ClubIdMetaKey, fmt.Sprintf("%d", in.SlotId))
	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err := client.RPCCallHash(_client, _client.ClubDelete, ctx, &in, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ClubRead(ctx context.Context, clubid int64, opts ...grpc.CallOption) (*pb.ClubDB, error) {
	var in pb.ClubRead_Req
	in.SlotId = club.Club2SlotId(clubid)
	in.ClubId = clubid
	md := metadata.Pairs(club.ClubIdMetaKey, fmt.Sprintf("%d", in.SlotId))
	ctx = metadata.NewOutgoingContext(ctx, md)
	out, err := client.RPCCallHash(_client, _client.ClubRead, ctx, &in, opts...)
	if err != nil {
		return nil, err
	} else {
		return out.ClubDB, nil
	}
}

func ClubJoin(ctx context.Context, clubid int64, roleid int64, opts ...grpc.CallOption) error {
	var in pb.ClubJoin_Req
	in.SlotId = club.Club2SlotId(clubid)
	in.ClubId = clubid
	in.RoleId = roleid
	md := metadata.Pairs(club.ClubIdMetaKey, fmt.Sprintf("%d", in.SlotId))
	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err := client.RPCCallHash(_client, _client.ClubJoin, ctx, &in, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ClubExit(ctx context.Context, clubid int64, roleid int64, opts ...grpc.CallOption) error {
	var in pb.ClubExit_Req
	in.SlotId = club.Club2SlotId(clubid)
	in.ClubId = clubid
	in.RoleId = roleid
	md := metadata.Pairs(club.ClubIdMetaKey, fmt.Sprintf("%d", in.SlotId))
	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err := client.RPCCallHash(_client, _client.ClubExit, ctx, &in, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ClubRoleInfo(ctx context.Context, clubid int64, roleid int64, opts ...grpc.CallOption) (*pb.ClubRoleDB, error) {
	var in pb.ClubRoleInfo_Req
	in.SlotId = club.Club2SlotId(clubid)
	in.ClubId = clubid
	in.RoleId = roleid
	md := metadata.Pairs(club.ClubIdMetaKey, fmt.Sprintf("%d", in.SlotId))
	ctx = metadata.NewOutgoingContext(ctx, md)
	out, err := client.RPCCallHash(_client, _client.ClubRoleInfo, ctx, &in, opts...)
	if err != nil {
		return nil, err
	} else {
		return out.ClubRoleDB, nil
	}
}
