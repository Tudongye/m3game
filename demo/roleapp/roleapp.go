package roleapp

import (
	"context"
	"m3game/config"
	"m3game/demo/onlineapp/onlinecli"
	"m3game/demo/proto"
	"m3game/demo/roleapp/rolecli"
	"m3game/demo/roleapp/roleser"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/db/redis"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	"m3game/plugins/router"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func newApp() *RoleApp {
	return &RoleApp{
		App: app.New(proto.RoleFuncID),
	}
}

var (
	_cfg AppCfg
)

type RoleApp struct {
	app.App
}

type AppCfg struct {
	PrePareTime int `mapstructure:"PrePareTime"`
}

func (c *AppCfg) CheckVaild() error {
	if c.PrePareTime == 0 {
		return errors.New("PrePareTime cant be 0")
	}
	return nil
}

func (a *RoleApp) Init(cfg map[string]interface{}) error {
	if err := mapstructure.Decode(cfg, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return err
	}
	return nil
}

func (d *RoleApp) Prepare(ctx context.Context) error {
	if err := rolecli.Init(config.GetAppID()); err != nil {
		return err
	}
	if err := onlinecli.Init(config.GetAppID()); err != nil {
		return err
	}
	return nil
}

func (d *RoleApp) Start(ctx context.Context) {
	log.Info("RoleApp PrepareTime %d", _cfg.PrePareTime)
	time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
	log.Info("RoleApp Ready")
	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// 插件检查
			if router.Get().Factory().CanDelete(router.Get()) {
				t.Stop()
				runtime.ShutDown("Router Delete")
			}
			continue
		}
	}
}

func (d *RoleApp) HealthCheck() bool {
	return true
}

func Run(ctx context.Context) error {
	runtime.Run(ctx, newApp(), []server.Server{roleser.New()})
	return nil
}
