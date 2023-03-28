package clubroleapp

import (
	"context"
	"m3game/demo/clubroleapp/clubroleser"
	"m3game/demo/proto"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/db/mongo"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	"m3game/plugins/router"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	_cfg AppCfg
)

func newApp() *ClubRoleApp {
	return &ClubRoleApp{
		App: app.New(proto.ClubRoleFuncID),
	}
}

type ClubRoleApp struct {
	app.App
}

type AppCfg struct {
	PrePareTime int `mapstructure:"PrePareTime" validate:"gt=0"`
}

func (a *ClubRoleApp) Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&_cfg); err != nil {
		return err
	}
	return nil
}

func (d *ClubRoleApp) Start(ctx context.Context) {
	log.Info("ClubRoleApp PrepareTime %d", _cfg.PrePareTime)
	time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
	log.Info("ClubRoleApp Ready")
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// 插件检查
			if router.Instance().Factory().CanUnload(router.Instance()) {
				runtime.ShutDown("Router Delete")
				return
			}
			continue
		}
	}
}

func (d *ClubRoleApp) HealthCheck() bool {
	return true
}

func Run(ctx context.Context) error {
	runtime.New().Run(ctx, newApp(), []server.Server{clubroleser.New()})
	return nil
}
