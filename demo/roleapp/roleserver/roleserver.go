package roleserver

import (
	"context"
	"errors"
	"fmt"
	"m3game/demo/dirapp/dirclient"
	"m3game/demo/loader"
	"m3game/demo/mapapp/mapclient"
	dpb "m3game/demo/proto/pb"
	"m3game/runtime/transport"
	"m3game/server/actor"
)

var (
	_map                 map[string]int32
	_err_actor_parsefail = errors.New("_err_actor_parsefail")
	_err_actor_created   = errors.New("_err_actor_created")
	_err_actor_readyed   = errors.New("_err_actor_created")
	_err_actor_notcreate = errors.New("_err_actor_notcreate")
	_err_actor_notready  = errors.New("_err_actor_notready")
	_err_actor_dberr     = errors.New("_err_actor_dberr")
)

func init() {
	_map = make(map[string]int32)
}

func CreateRoleSer() *RoleSer {
	return &RoleSer{
		Server: actor.CreateServer("RoleSer", RoleCreater),
	}
}

type RoleSer struct {
	*actor.Server
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
	out.RoleID = role.roleid
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
	if rsp, err := dirclient.DirClient().Hello(ctx, role.Name()); err != nil {
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
	if _, l, err := mapclient.MapClient().Move(ctx, role.Name(), in.Distance); err != nil {
		return nil, fmt.Errorf("Move Fail err:%s", err.Error())
	} else {
		out.Location = l
		locationcfgloader := loader.GetLocationCfgLoader()
		if locationcfgloader == nil {
			return out, fmt.Errorf("LocationCfg Err")
		}
		out.LocateName = locationcfgloader.GetNameByDistance(l)
	}
	return out, nil
}

func (s *RoleSer) TransportRegister() func(*transport.Transport) error {
	return func(t *transport.Transport) error {
		dpb.RegisterRoleSerServer(t.GrpcSer(), s)
		return nil
	}
}
