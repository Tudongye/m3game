package roleser

import (
	"context"
	"fmt"
	"m3game/demo/proto/pb"
	"m3game/demo/roleapp/role"
	"m3game/meta"
	"m3game/plugins/log"
	"m3game/runtime/rpc"
	mactor "m3game/runtime/server/actor"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	_err_actor_parsefail = errors.New("_err_actor_parsefail")
	_err_actor_readyed   = errors.New("_err_actor_readyed")
	_err_actor_notready  = errors.New("_err_actor_notready")
)

var (
	_ser *RoleSer
)

func init() {
	if err := rpc.InjectionRPC(pb.File_role_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("InjectionRPC Role %s", err.Error()))
	}
}

func New() *RoleSer {
	if _ser != nil {
		return _ser
	}
	_ser = &RoleSer{
		Server: mactor.New("RoleSer", role.RoleCreater),
	}
	return _ser
}

func RoleLogic(ctx context.Context, f func(role *role.Role) error) error {
	role := role.ConvertRole(ctx)
	if role == nil || !role.Ready() {
		return _err_actor_notready
	}
	return f(role)
}

type RoleSer struct {
	*mactor.Server
	pb.UnimplementedRoleSerServer
}

func (s *RoleSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterRoleSerServer(t, s)
		return nil
	}
}

func (d *RoleSer) RoleLogin(ctx context.Context, in *pb.RoleLogin_Req) (*pb.RoleLogin_Rsp, error) {
	out := new(pb.RoleLogin_Rsp)
	log.Info("Login")
	role := role.ConvertRole(ctx)
	if role == nil {
		return out, _err_actor_parsefail
	}
	if role.Ready() {
		log.Debug(role.ActorID())
		return out, _err_actor_readyed
	}
	if err := role.Login(ctx); err != nil {
		return out, err
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vlist, ok := md[string(meta.M3RouteSrcApp)]; ok {
			if len(vlist) > 0 {
				log.Debug("Actor %s at Gate %s", role.ActorID(), vlist[0])
				role.SetGate(vlist[0])
			}
		}
	}
	out.RoleDB = role.DB()
	return out, nil
}

func (d *RoleSer) RoleGetInfo(ctx context.Context, in *pb.RoleGetInfo_Req) (*pb.RoleGetInfo_Rsp, error) {
	out := new(pb.RoleGetInfo_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		out.RoleDB = role.DB()
		return nil
	})
}

func (d *RoleSer) RoleModifyName(ctx context.Context, in *pb.RoleModifyName_Req) (*pb.RoleModifyName_Rsp, error) {
	out := new(pb.RoleModifyName_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		role.ModifyName(in.NewName)
		out.Name = role.DB().GetName()
		return nil
	})
}

func (d *RoleSer) RolePowerUp(ctx context.Context, in *pb.RolePowerUp_Req) (*pb.RolePowerUp_Rsp, error) {
	out := new(pb.RolePowerUp_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		role.PowerUp(in.PowerUp)
		out.Power = role.DB().GetPower()
		return nil
	})
}

func (d *RoleSer) RoleKick(ctx context.Context, in *pb.RoleKick_Req) (*pb.RoleKick_Rsp, error) {
	out := new(pb.RoleKick_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		return role.Kick(ctx)
	})
}
