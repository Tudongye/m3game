package onlineser

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
	_cfg OnlineSerCfg
	_ser *OnlineSer
)

func init() {
	if err := rpc.InjectionRPC(pb.File_online_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC OnlineSer %s", err.Error())
	}
}

type OnlineSerCfg struct {
	CachePoolSize   int `mapstructure:"CachePoolSize" validate:"gt=0"`
	AppAliveTimeOut int `mapstructure:"AppAliveTimeOut" validate:"gt=0"`
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

func New() *OnlineSer {
	if _ser != nil {
		return _ser
	}
	_ser = &OnlineSer{
		Server: multi.New("OnlineSer"),
		pool:   newPool(),
	}
	return _ser
}

type OnlineSer struct {
	*multi.Server
	pb.UnimplementedOnlineSerServer
	pool *OnlinePool
}

func (s *OnlineSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterOnlineSerServer(t, s)
		return nil
	}
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
