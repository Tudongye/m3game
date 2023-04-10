package multiapp

import (
	"context"
	"m3game/example/multiapp/multiser"
	"m3game/example/proto"
	_ "m3game/plugins/broker/nats"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	"m3game/plugins/router"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
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

var (
	_cfg AppCfg
)

func newApp() *MultiApp {
	return &MultiApp{
		App: app.New(proto.MultiAppFuncID),
	}
}

type MultiApp struct {
	app.App
}

type AppCfg struct {
	PrePareTime int `mapstructure:"PrePareTime" validate:"gt=0"`
}

func (a *MultiApp) Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&_cfg); err != nil {
		return err
	}
	return nil
}

func (d *MultiApp) Start(ctx context.Context) {
	log.Info("MultiApp PrepareTime %d", _cfg.PrePareTime)
	time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
	log.Info("MultiApp Ready")
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

func (d *MultiApp) Alive(app string, svc string) bool {
	return true
}

func Run(ctx context.Context) error {
	runtime.New().Run(ctx, newApp(), []server.Server{multiser.New()})
	return nil
}
