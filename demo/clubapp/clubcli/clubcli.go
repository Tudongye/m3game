package clubcli

import (
	"context"
	"fmt"
	"m3game/demo/clubapp/club"
	"m3game/demo/proto"
	"m3game/demo/proto/pb"
	"m3game/meta"
	"m3game/plugins/log"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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

func New(srcapp meta.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := meta.GenDstRouteSvc(srcapp, proto.ClubFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	target := fmt.Sprintf("router://%s", _client.DstSvc().String())
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(client.ClientInterceptors()...)),
		grpc.WithTimeout(time.Second*10),
	)
	if _client.conn, err = grpc.Dial(target, opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.ClubSerClient = pb.NewClubSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.ClubSerClient
	conn *grpc.ClientConn
}

func Conn() *grpc.ClientConn {
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
