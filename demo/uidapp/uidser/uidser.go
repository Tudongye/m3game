package uidser

import (
	"context"
	"m3game/demo/proto/pb"
	"m3game/plugins/log"
	"m3game/runtime/rpc"
	"m3game/runtime/server/multi"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_cfg UidSerCfg
	_ser *UidSer
)

func init() {
	if err := rpc.InjectionRPC(pb.File_uid_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC UidSer %s", err.Error())
	}
}

type UidSerCfg struct {
	CachePoolSize int `mapstructure:"CachePoolSize" validate:"gt=0"`
}

func Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errors.Wrapf(err, "App Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&_cfg); err != nil {
		return err
	}
	return nil
}

func New() *UidSer {
	if _ser != nil {
		return _ser
	}
	_ser = &UidSer{
		Server: multi.New("UidSer"),
		pool:   newPool(),
	}
	return _ser
}

type UidSer struct {
	*multi.Server
	pb.UnimplementedUidSerServer
	pool *UidPool
}

func (s *UidSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterUidSerServer(t, s)
		return nil
	}
}

func (d *UidSer) AllocRoleId(ctx context.Context, in *pb.AllocRoleId_Req) (*pb.AllocRoleId_Rsp, error) {
	out := new(pb.AllocRoleId_Rsp)
	if roleid, err := _uidpool.AllocRoleId(ctx, in.OpenId); err != nil {
		return out, err
	} else {
		out.RoleId = roleid
	}
	return out, nil
}

func (d *UidSer) AllocClubId(ctx context.Context, in *pb.AllocClubId_Req) (*pb.AllocClubId_Rsp, error) {
	out := new(pb.AllocClubId_Rsp)
	if clubid, err := _uidpool.AllocClubId(ctx, in.RoleId); err != nil {
		return out, err
	} else {
		out.ClubId = clubid
	}
	return out, nil
}
