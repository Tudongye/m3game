package clubdser

import (
	"context"
	"m3game/demo/clubapp/club"
	"m3game/demo/proto/pb"
	"m3game/plugins/log"
	"m3game/runtime/rpc"
	"m3game/runtime/server/async"

	"google.golang.org/grpc"
)

func init() {
	if err := rpc.InjectionRPC(pb.File_club_proto.Services().Get(1)); err != nil {
		log.Fatal("InjectionRPC ClubDSer %s", err.Error())
	}
}

func New() *ClubDSer {
	return &ClubDSer{
		Server: async.New("ClubDSer"),
	}
}

type ClubDSer struct {
	*async.Server
	pb.UnimplementedClubDaemonSerServer
}

func (d *ClubDSer) ClubKick(ctx context.Context, in *pb.ClubKick_Req) (*pb.ClubKick_Rsp, error) {
	out := new(pb.ClubKick_Rsp)
	log.Info("Kick")
	if _, err := club.LeaseMeta().RecvKickLease(ctx, in.LeaseId); err != nil {
		return out, err
	}
	return out, nil
}

func (d *ClubDSer) ClubList(ctx context.Context, in *pb.ClubList_Req) (*pb.ClubList_Rsp, error) {
	out := new(pb.ClubList_Rsp)
	return out, nil
}

func (s *ClubDSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterClubDaemonSerServer(t, s)
		return nil
	}
}
