package onlineser

import (
	"context"
	"fmt"
	"m3game/demo/proto/pb"
	"m3game/runtime/rpc"
	"m3game/runtime/server/multi"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_cfg OnlineSerCfg
)

func init() {
	if err := rpc.RegisterRPCSvc(pb.File_online_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc OnlineSer %s", err.Error()))
	}
}

type OnlineSerCfg struct {
	CachePoolSize   int `mapstructure:"CachePoolSize"`
	AppAliveTimeOut int `mapstructure:"AppAliveTimeOut"`
}

func (c OnlineSerCfg) CheckVaild() error {
	if c.CachePoolSize == 0 {
		return errors.New("CachePoolSize cant be 0")
	}
	if c.AppAliveTimeOut == 0 {
		return errors.New("AppAliveTimeOut cant be 0")
	}
	return nil
}

func Init(cfg map[string]interface{}) error {
	if err := mapstructure.Decode(cfg, &_cfg); err != nil {
		return errors.Wrapf(err, "App Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return err
	}
	return nil
}

func New() *OnlineSer {
	return &OnlineSer{
		Server: multi.New("OnlineSer"),
	}
}

type OnlineSer struct {
	*multi.Server
	pb.UnimplementedOnlineSerServer
}

func (d *OnlineSer) OnlineCreate(ctx context.Context, in *pb.OnlineCreate_Req) (*pb.OnlineCreate_Rsp, error) {
	out := new(pb.OnlineCreate_Rsp)
	if err := _onlinepool.OnlineCreate(in.RoleId, in.AppId); err != nil {
		return out, err
	} else {
		return out, nil
	}
}
func (d *OnlineSer) OnlineDelete(ctx context.Context, in *pb.OnlineDelete_Req) (*pb.OnlineDelete_Rsp, error) {
	out := new(pb.OnlineDelete_Rsp)
	if err := _onlinepool.OnlineDelete(in.RoleId, in.AppId); err != nil {
		return out, err
	} else {
		return out, nil
	}
}
func (d *OnlineSer) OnlineRead(ctx context.Context, in *pb.OnlineRead_Req) (*pb.OnlineRead_Rsp, error) {
	out := new(pb.OnlineRead_Rsp)
	if appid, err := _onlinepool.OnlineRead(in.RoleId); err != nil {
		return out, err
	} else {
		out.AppId = appid
		return out, nil
	}
}

func (s *OnlineSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterOnlineSerServer(t, s)
		return nil
	}
}
