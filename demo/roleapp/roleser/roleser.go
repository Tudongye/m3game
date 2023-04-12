package roleser

import (
	"context"
	"fmt"
	"m3game/demo/clubapp/clubcli"
	"m3game/demo/proto/pb"
	"m3game/demo/roleapp/role"
	"m3game/demo/uidapp/uidcli"
	"m3game/meta"
	"m3game/plugins/log"
	"m3game/runtime/app"
	"m3game/runtime/rpc"
	mactor "m3game/runtime/server/actor"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
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
		log.Fatal("InjectionRPC Role %s", err.Error())
	}
}

func New() *RoleSer {
	if _ser != nil {
		return _ser
	}
	_ser = &RoleSer{
		Server: mactor.New("RoleSer", role.RoleCreater, role.RoleIdMetaKey),
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
	cfg Config
}

type Config struct {
	LuaFile string `mapstructure:"LuaFile" validate:"required"`
}

func (s *RoleSer) Init(c map[string]interface{}, app app.App) error {
	if err := s.Server.Init(c, app); err != nil {
		return err
	}
	if err := mapstructure.Decode(c, &s.cfg); err != nil {
		return fmt.Errorf("Decode RoleSer Config %s", err.Error())
	}
	validate := validator.New()
	if err := validate.Struct(&s.cfg); err != nil {
		return fmt.Errorf("Decode RoleSer Config %s", err.Error())
	}
	return nil
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
		return out, _err_actor_readyed
	}

	if ret, err := LuaRoleHook(d.cfg.LuaFile, role.ActorID()); err != nil {
		return out, err
	} else if !ret {
		return out, errors.New("")
	}

	if err := role.Login(ctx); err != nil {
		return out, err
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vlist, ok := md[string(meta.M3RouteSrcApp)]; ok {
			if len(vlist) > 0 {
				role.SetGate(vlist[0])
			}
		}
	}
	out.RoleDB = role.Detail()
	return out, nil
}

func (d *RoleSer) RoleGetInfo(ctx context.Context, in *pb.RoleGetInfo_Req) (*pb.RoleGetInfo_Rsp, error) {
	out := new(pb.RoleGetInfo_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		if in.Detail {
			out.RoleDB = role.Detail()
		} else {
			out.RoleDB = role.Brief()
		}
		if clubid := role.GetClubId(); clubid != 0 {
			if clubroledb, err := clubcli.ClubRoleInfo(ctx, clubid, role.DB().RoleId); err != nil {
				return err
			} else {
				out.ClubRoleDB = clubroledb
			}
		}
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

func (d *RoleSer) RoleGetClubInfo(ctx context.Context, in *pb.RoleGetClubInfo_Req) (*pb.RoleGetClubInfo_Rsp, error) {
	out := new(pb.RoleGetClubInfo_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		if clubid := role.GetClubId(); clubid == 0 {
			return nil
		} else {
			if clubdb, err := clubcli.ClubRead(ctx, clubid); err != nil {
				return err
			} else {
				out.ClubDB = clubdb
			}
		}
		return nil
	})
}

func (d *RoleSer) RoleGetClubList(ctx context.Context, in *pb.RoleGetClubList_Req) (*pb.RoleGetClubList_Rsp, error) {
	out := new(pb.RoleGetClubList_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		return nil
	})
}

func (d *RoleSer) RoleCreateClub(ctx context.Context, in *pb.RoleCreateClub_Req) (*pb.RoleCreateClub_Rsp, error) {
	out := new(pb.RoleCreateClub_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		if role.GetClubId() != 0 {
			return errors.New("Role Has in Club")
		}
		if clubid, err := uidcli.AllocClubId(ctx, role.DB().RoleId); err != nil {
			return err
		} else {
			if err := clubcli.ClubCreate(ctx, clubid, role.DB().RoleId); err != nil {
				return err
			}
			out.ClubId = clubid
			role.SetClubId(clubid)
		}
		return nil
	})
}

func (d *RoleSer) RoleJoinClub(ctx context.Context, in *pb.RoleJoinClub_Req) (*pb.RoleJoinClub_Rsp, error) {
	out := new(pb.RoleJoinClub_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		if role.GetClubId() != 0 {
			return errors.New("Role Has in Club")
		}
		if err := clubcli.ClubJoin(ctx, in.ClubId, role.DB().RoleId); err != nil {
			return err
		}
		return nil
	})
}

func (d *RoleSer) RoleExitClub(ctx context.Context, in *pb.RoleExitClub_Req) (*pb.RoleExitClub_Rsp, error) {
	out := new(pb.RoleExitClub_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		if role.GetClubId() == 0 {
			return errors.New("Role not in Club")
		}
		if err := clubcli.ClubExit(ctx, role.GetClubId(), role.DB().RoleId); err != nil {
			return err
		}
		role.SetClubId(0)
		return nil
	})
}

func (d *RoleSer) RoleCancelClub(ctx context.Context, in *pb.RoleCancelClub_Req) (*pb.RoleCancelClub_Rsp, error) {
	out := new(pb.RoleCancelClub_Rsp)
	return out, RoleLogic(ctx, func(role *role.Role) error {
		if role.GetClubId() == 0 {
			return errors.New("Role not in Club")
		}
		if err := clubcli.ClubDelete(ctx, role.GetClubId(), role.DB().RoleId); err != nil {
			return err
		}
		role.SetClubId(0)
		return nil
	})
}
