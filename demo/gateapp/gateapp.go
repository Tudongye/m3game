package gateapp

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/demo/gateapp/gateser"
	"m3game/demo/proto"
	"m3game/demo/roleapp/rolecli"
	"m3game/demo/uidapp/uidcli"
	"m3game/meta"
	"m3game/meta/metapb"
	_ "m3game/plugins/broker/nats"
	"m3game/plugins/gate"
	_ "m3game/plugins/gate/grpcgate"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/metric/prometheus"
	"m3game/plugins/router"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/rpc"
	"m3game/runtime/server"
	"m3game/util"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_cfg AppCfg
)

func newApp() *GateApp {
	return &GateApp{
		App: app.New(proto.GateFuncID),
	}
}

type GateApp struct {
	app.App
}

type AppCfg struct {
	PrePareTime int `mapstructure:"PrePareTime"`
}

func (c *AppCfg) checkValid() error {
	if err := util.InEqualInt(c.PrePareTime, 0, "PrePareTime"); err != nil {
		return err
	}
	return nil
}

func (a *GateApp) Init(cfg map[string]interface{}) error {
	if err := mapstructure.Decode(cfg, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	if err := _cfg.checkValid(); err != nil {
		return err
	}
	return nil
}

func (d *GateApp) Prepare(ctx context.Context) error {
	if err := rolecli.Init(config.GetAppID(), grpc.WithCodec(&gate.GateCodec{})); err != nil {
		return err
	}
	if err := uidcli.Init(config.GetAppID()); err != nil {
		return err
	}
	gate.SetReciver(d)
	return nil
}

func (d *GateApp) Start(ctx context.Context) {
	log.Info("GateApp PrepareTime %d", _cfg.PrePareTime)
	time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
	log.Info("GateApp Ready")
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// 插件检查
			if router.Get().Factory().CanDelete(router.Get()) {
				runtime.ShutDown("Router Delete")
				return
			}
			if gate.Get().Factory().CanDelete(gate.Get()) {
				runtime.ShutDown("Gate Delete")
				return
			}
			continue
		}
	}
}

func (d *GateApp) LogicCall(roleid string, in *metapb.CSMsg) (*metapb.CSMsg, error) {
	if !rpc.IsRPCClientMethod(in.Method) {
		return nil, fmt.Errorf("Method %s invaild", in.Method)
	}
	// 路由参数
	in.Metas = append(in.Metas, &metapb.Meta{Key: meta.M3RouteType.String(), Value: meta.RouteTypeHash.String()})
	in.Metas = append(in.Metas, &metapb.Meta{Key: meta.M3RouteHashKey.String(), Value: roleid})

	var out *metapb.CSMsg
	var err error
	if strings.HasPrefix(in.Method, "/proto.RoleSer") {
		out, err = gate.CallGrpcCli(context.Background(), rolecli.Conn(), in)
	}
	if err != nil {
		return out, err
	} else if out != nil {
		out.Metas = in.Metas
		return out, nil
	}
	return nil, fmt.Errorf("Unknow Method %s", in.Method)
}

func (d *GateApp) AuthCall(req *metapb.AuthReq) (*metapb.AuthRsp, error) {
	log.Debug("AuthCall...")
	openid := fmt.Sprintf("OpenId-%s", req.Token)
	if roleid, err := uidcli.AllocRoleId(context.Background(), openid); err != nil {
		return nil, err
	} else {
		return &metapb.AuthRsp{
			PlayerID: roleid,
		}, nil
	}
}

func (d *GateApp) HealthCheck() bool {
	return true
}

func Run(ctx context.Context) error {
	runtime.Run(ctx, newApp(), []server.Server{gateser.New()})
	return nil
}
