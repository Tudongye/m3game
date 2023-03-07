package roleserver

import (
	"context"
	"fmt"
	"m3game/demo/dirapp/dirclient"
	"m3game/demo/loader"
	"m3game/demo/mapapp/mapclient"
	dpb "m3game/demo/proto/pb"
	"m3game/demo/roleapp/rolechclient"
	"m3game/demo/roleapp/rolechserver"
	"m3game/runtime/resource"
	"m3game/runtime/rpc"
	"m3game/runtime/server/actor"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
)

var (
	_err_actor_parsefail = errors.New("_err_actor_parsefail")
	_err_actor_created   = errors.New("_err_actor_created")
	_err_actor_readyed   = errors.New("_err_actor_created")
	_err_actor_notcreate = errors.New("_err_actor_notcreate")
	_err_actor_notready  = errors.New("_err_actor_notready")
	_err_actor_dberr     = errors.New("_err_actor_dberr")
)

func init() {
	if err := rpc.RegisterRPCSvc(dpb.File_role_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc Role %s", err.Error()))
	}
}

func New() *RoleSer {
	return &RoleSer{
		Server: actor.New("RoleSer", roleCreater),
	}
}

type RoleSer struct {
	*actor.Server
	dpb.UnimplementedRoleSerServer
}

func (s *RoleSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		dpb.RegisterRoleSerServer(t, s)
		return nil
	}
}

func (d *RoleSer) Register(ctx context.Context, in *dpb.Register_Req) (*dpb.Register_Rsp, error) {
	out := new(dpb.Register_Rsp)
	role := ParseRoleActor(ctx)
	if role == nil {
		return out, _err_actor_parsefail
	}
	if err := role.Register(in.Name); err != nil {
		return out, err
	}
	out.RoleID = role.RoleID()
	return out, nil
}

func (d *RoleSer) Login(ctx context.Context, in *dpb.Login_Req) (*dpb.Login_Rsp, error) {
	out := new(dpb.Login_Rsp)
	role := ParseRoleActor(ctx)
	if role == nil {
		return out, _err_actor_parsefail
	}
	if role.ready {
		return out, _err_actor_readyed
	}
	if err := role.Login(); err != nil {
		return out, err
	}
	out.Name = role.Name()
	if rsp, err := dirclient.Hello(ctx, role.Name()); err != nil {
		return out, fmt.Errorf("Hello Fail err:%s", err.Error())
	} else {
		out.Tips = rsp
	}
	return out, nil
}

func (d *RoleSer) ModifyName(ctx context.Context, in *dpb.ModifyName_Req) (*dpb.ModifyName_Rsp, error) {
	out := new(dpb.ModifyName_Rsp)
	role := ParseRoleActor(ctx)
	if role == nil || !role.ready {
		return out, _err_actor_notready
	}
	role.ModifyName(in.NewName)
	out.Name = role.Name()
	return out, nil
}

func (d *RoleSer) GetName(ctx context.Context, in *dpb.GetName_Req) (*dpb.GetName_Rsp, error) {
	out := new(dpb.GetName_Rsp)
	role := ParseRoleActor(ctx)
	if role == nil || !role.ready {
		return out, _err_actor_notready
	}
	out.Name = role.Name()
	return out, nil
}

func (d *RoleSer) MoveRole(ctx context.Context, in *dpb.MoveRole_Req) (*dpb.MoveRole_Rsp, error) {
	out := new(dpb.MoveRole_Rsp)
	role := ParseRoleActor(ctx)
	if role == nil || !role.ready {
		return out, _err_actor_notready
	}
	if _, l, err := mapclient.Move(ctx, role.Name(), in.Distance); err != nil {
		return nil, fmt.Errorf("Move Fail err:%s", err.Error())
	} else {
		out.Location = l
		locationcfgloader := resource.GetLoader[*loader.LocationCfgLoader](ctx)
		if locationcfgloader == nil {
			return out, fmt.Errorf("LocationCfg Err")
		}
		out.LocateName = locationcfgloader.GetNameByDistance(l)
		role.ModifyLocation(out.LocateName, out.Location)
	}
	return out, nil
}

func (d *RoleSer) PostChannel(ctx context.Context, in *dpb.PostChannel_Req) (*dpb.PostChannel_Rsp, error) {
	out := new(dpb.PostChannel_Rsp)
	role := ParseRoleActor(ctx)
	if role == nil || !role.ready {
		return out, _err_actor_notready
	}
	if err := rolechclient.TransChannel(ctx, &dpb.ChannelMsg{Name: role.Name(), Content: in.Content}); err != nil {
		return out, err
	} else {
		return out, nil
	}
}

func (d *RoleSer) PullChannel(ctx context.Context, in *dpb.PullChannel_Req) (*dpb.PullChannel_Rsp, error) {
	out := new(dpb.PullChannel_Rsp)
	role := ParseRoleActor(ctx)
	if role == nil || !role.ready {
		return out, _err_actor_notready
	}
	out.Msgs = rolechserver.GetMsg()
	return out, nil
}
