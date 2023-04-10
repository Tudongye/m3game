package roleapp

import (
	"context"
	"m3game/config"
	"m3game/demo/clubapp/clubcli"
	"m3game/demo/clubroleapp/clubrolecli"
	"m3game/demo/onlineapp/onlinecli"
	"m3game/demo/proto"
	"m3game/demo/roleapp/rolecli"
	"m3game/demo/roleapp/roleser"
	"m3game/demo/uidapp/uidcli"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/db/mongo"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/metric/prometheus"
	"m3game/plugins/router"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	_ "m3game/plugins/trace/jaeger"
	_ "m3game/plugins/transport/natstrans"
	_ "m3game/plugins/transport/tcptrans"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func newApp() *RoleApp {
	return &RoleApp{
		App: app.New(proto.RoleFuncID),
	}
}

type RoleApp struct {
	app.App
	cfg AppCfg
}

type AppCfg struct {
	PrePareTime int `mapstructure:"PrePareTime" validate:"gt=0"`
}

func (a *RoleApp) Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &a.cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&a.cfg); err != nil {
		return err
	}

	return nil
}

func (d *RoleApp) Prepare(ctx context.Context) error {
	if _, err := rolecli.New(config.GetAppID()); err != nil {
		return err
	} else if _, err := onlinecli.New(config.GetAppID()); err != nil {
		return err
	} else if _, err := uidcli.New(config.GetAppID()); err != nil {
		return err
	} else if _, err := clubrolecli.New(config.GetAppID()); err != nil {
		return err
	} else if _, err := clubcli.New(config.GetAppID()); err != nil {
		return err
	}
	return nil
}

func (a *RoleApp) Start(ctx context.Context) {
	log.Info("RoleApp PrepareTime %d", a.cfg.PrePareTime)
	time.Sleep(time.Duration(a.cfg.PrePareTime) * time.Second)
	log.Info("RoleApp Ready")
	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// 插件检查
			if router.Instance().Factory().CanUnload(router.Instance()) {
				t.Stop()
				runtime.ShutDown("Router Delete")
			}
			continue
		}
	}
}

func (a *RoleApp) Alive(app string, svc string) bool {
	return true
}

func Run(ctx context.Context) error {
	return runtime.New().Run(ctx, newApp(), []server.Server{roleser.New()})
}
