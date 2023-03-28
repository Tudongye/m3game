package clubser

import (
	"context"
	"fmt"
	"m3game/demo/clubapp/club"
	"m3game/demo/proto/pb"
	"m3game/plugins/log"
	"m3game/runtime/rpc"
	"m3game/runtime/server/actor"

	"google.golang.org/grpc"
)

var (
	_ser *ClubSer
)

var (
	Err_SlotNotFind = fmt.Errorf("Err_SlotNotFind")
	Err_ClubNotFind = fmt.Errorf("Err_ClubNotFind")
)

func init() {
	if err := rpc.InjectionRPC(pb.File_club_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC Club %s", err.Error())
	}
}

func New() *ClubSer {
	if _ser != nil {
		return _ser
	}
	_ser = &ClubSer{
		Server: actor.New("ClubSer", club.SlotCreater, club.ClubIdMetaKey),
	}
	return _ser
}

func Ser() *ClubSer {
	return _ser
}

func ClubLogic(ctx context.Context, clubid int64, f func(role *club.Club) error) error {
	slot := club.ConvertSlot(ctx)
	if slot == nil {
		return Err_SlotNotFind
	}
	if club := slot.Club(clubid); club == nil {
		return Err_ClubNotFind
	} else {
		return f(club)
	}
}

type ClubSer struct {
	*actor.Server
	pb.UnimplementedClubSerServer
}

func (s *ClubSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterClubSerServer(t, s)
		return nil
	}
}

func (d *ClubSer) ClubCreate(ctx context.Context, in *pb.ClubCreate_Req) (*pb.ClubCreate_Rsp, error) {
	out := new(pb.ClubCreate_Rsp)
	clubslot := club.ConvertSlot(ctx)
	if clubslot == nil {
		return out, Err_SlotNotFind
	}
	if err := clubslot.CreateClub(ctx, in.ClubId, in.RoleId); err != nil {
		return out, err
	}
	return out, nil
}

func (d *ClubSer) ClubDelete(ctx context.Context, in *pb.ClubDelete_Req) (*pb.ClubDelete_Rsp, error) {
	out := new(pb.ClubDelete_Rsp)
	clubslot := club.ConvertSlot(ctx)
	if clubslot == nil {
		return out, Err_SlotNotFind
	}
	if err := clubslot.DeleteClub(ctx, in.ClubId, in.RoleId); err != nil {
		return out, err
	}
	return out, nil
}

func (d *ClubSer) ClubRead(ctx context.Context, in *pb.ClubRead_Req) (*pb.ClubRead_Rsp, error) {
	out := new(pb.ClubRead_Rsp)
	return out, ClubLogic(ctx, in.ClubId, func(club *club.Club) error {
		out.ClubDB = club.Obj()
		return nil
	})
}

func (d *ClubSer) ClubJoin(ctx context.Context, in *pb.ClubJoin_Req) (*pb.ClubJoin_Rsp, error) {
	out := new(pb.ClubJoin_Rsp)
	return out, ClubLogic(ctx, in.ClubId, func(club *club.Club) error {
		return club.Join(ctx, in.RoleId)
	})
}

func (d *ClubSer) ClubExit(ctx context.Context, in *pb.ClubExit_Req) (*pb.ClubExit_Rsp, error) {
	out := new(pb.ClubExit_Rsp)
	return out, ClubLogic(ctx, in.ClubId, func(club *club.Club) error {
		return club.Exit(ctx, in.RoleId)
	})
}

func (d *ClubSer) ClubRoleInfo(ctx context.Context, in *pb.ClubRoleInfo_Req) (*pb.ClubRoleInfo_Rsp, error) {
	out := new(pb.ClubRoleInfo_Rsp)
	return out, ClubLogic(ctx, in.ClubId, func(club *club.Club) error {
		clubrole := club.ClubRole(in.RoleId)
		if clubrole != nil {
			out.ClubRoleDB = clubrole.Obj()
		}
		return nil
	})
}
