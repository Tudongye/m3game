package uidser

import (
	"context"
	"fmt"
	"m3game/demo/proto/pb"
	"m3game/runtime/rpc"
	"m3game/runtime/server/multi"
	"m3game/util"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_cfg UidSerCfg
)

func init() {
	if err := rpc.InjectionRPC(pb.File_uid_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("InjectionRPC UidSer %s", err.Error()))
	}
}

type UidSerCfg struct {
	CachePoolSize int `mapstructure:"CachePoolSize"`
}

func (c UidSerCfg) checkValid() error {
	if err := util.InEqualInt(c.CachePoolSize, 0, "CachePoolSize"); err != nil {
		return err
	}
	return nil
}

func Init(cfg map[string]interface{}) error {
	if err := mapstructure.Decode(cfg, &_cfg); err != nil {
		return errors.Wrapf(err, "App Decode Cfg")
	}
	if err := _cfg.checkValid(); err != nil {
		return err
	}
	return nil
}

func New() *UidSer {
	return &UidSer{
		Server: multi.New("UidSer"),
	}
}

type UidSer struct {
	*multi.Server
	pb.UnimplementedUidSerServer
}

func (d *UidSer) AllocRoleId(ctx context.Context, in *pb.AllocRoleId_Req) (*pb.AllocRoleId_Rsp, error) {
	out := new(pb.AllocRoleId_Rsp)
	if roleid, err := _uidpool.AllocRoleId(in.OpenId); err != nil {
		return out, err
	} else {
		out.RoleId = roleid
	}
	return out, nil
}

func (d *UidSer) AllocClubId(ctx context.Context, in *pb.AllocClubId_Req) (*pb.AllocClubId_Rsp, error) {
	out := new(pb.AllocClubId_Rsp)
	if clubid, err := _uidpool.AllocClubId(in.RoleId); err != nil {
		return out, err
	} else {
		out.ClubId = clubid
	}
	return out, nil
}

func (s *UidSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterUidSerServer(t, s)
		return nil
	}
}
