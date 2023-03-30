package clubroleser

import (
	"context"
	"m3game/demo/proto/pb"
	"m3game/meta/errs"
	"m3game/plugins/db"
	"m3game/plugins/log"
	"m3game/runtime/rpc"
	"m3game/runtime/server/multi"

	"google.golang.org/grpc"
)

func init() {
	if err := rpc.InjectionRPC(pb.File_clubrole_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC MultiSer %s", err.Error())
	}
}

func New() *ClubRoleSer {
	return &ClubRoleSer{
		Server: multi.New("ClubRoleSer"),
	}
}

type ClubRoleSer struct {
	*multi.Server
	pb.UnimplementedClubRoleSerServer
}

func (s *ClubRoleSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterClubRoleSerServer(t, s)
		return nil
	}
}

func (d *ClubRoleSer) ClubRoleRead(ctx context.Context, in *pb.ClubRoleRead_Req) (*pb.ClubRoleRead_Rsp, error) {
	out := new(pb.ClubRoleRead_Rsp)
	w := _clubrolewrapermeta.New(in.RoleId)
	dbplugin := db.Instance()
	if err := w.Read(ctx, dbplugin); err != nil {
		if errs.DBKeyNotFound.Is(err) {
			return out, nil
		}
		return out, err
	}
	out.ClubId = w.Obj().ClubId
	return out, nil
}
