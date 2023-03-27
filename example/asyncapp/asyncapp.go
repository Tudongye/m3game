package asyncapp

import (
	"context"
	"m3game/example/asyncapp/asyncser"
	"m3game/example/proto"
	_ "m3game/plugins/broker/nats"
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

func newApp() *AsyncApp {
	return &AsyncApp{
		App: app.New(proto.AsyncAppFuncID),
	}
}

type AsyncApp struct {
	app.App
}

type AppCfg struct {
	PrePareTime int `mapstructure:"PrePareTime" validate:"gt=0"`
}

func (a *AsyncApp) Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&_cfg); err != nil {
		return err
	}
	return nil
}
func (d *AsyncApp) HealthCheck() bool {
	return true
}

func (d *AsyncApp) Start(ctx context.Context) {
	log.Info("AsyncApp PrepareTime %d", _cfg.PrePareTime)
	time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
	log.Info("AsyncApp Ready")
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

func Run(ctx context.Context) error {
	runtime.New().Run(ctx, newApp(), []server.Server{asyncser.New()})
	return nil
}
