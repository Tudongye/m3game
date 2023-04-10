package actorapp

import (
	"context"
	"m3game/config"
	"m3game/example/actorapp/actorregcli"
	"m3game/example/actorapp/actorregser"
	"m3game/example/actorapp/actorser"
	"m3game/example/asyncapp/asynccli"
	"m3game/example/gateapp/gatecli"
	"m3game/example/loader"
	"m3game/example/proto"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/db/cache"
	"m3game/plugins/lease"
	_ "m3game/plugins/lease/etcd"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	"m3game/plugins/router"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	_ "m3game/plugins/transport/tcptrans"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func newApp() *ActorApp {
	return &ActorApp{
		App: app.New(proto.ActorAppFuncID),
	}
}

var (
	_cfg AppCfg
)

type ActorApp struct {
	app.App
}

type AppCfg struct {
	PrePareTime int `mapstructure:"PrePareTime" validate:"gt=0"`
}

func (a *ActorApp) Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&_cfg); err != nil {
		return err
	}
	return nil
}

func (d *ActorApp) Prepare(ctx context.Context) error {
	if err := asynccli.Init(config.GetAppID()); err != nil {
		return err
	} else if err := gatecli.Init(config.GetAppID()); err != nil {
		return err
	} else if err := actorregcli.Init(config.GetAppID()); err != nil {
		return err
	}
	lease.SetReciver(d)
	return nil
}

func (d *ActorApp) Start(ctx context.Context) {
	log.Info("ActorApp PrepareTime %d", _cfg.PrePareTime)
	time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
	log.Info("ActorApp Ready")
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
			if lease.Instance().Factory().CanUnload(lease.Instance()) {
				t.Stop()
				runtime.ShutDown("Lease Delete")
			}
			continue
		}
	}
}

func (d *ActorApp) Alive(app string, svc string) bool {
	return true
}

func (d *ActorApp) PreExitLease(ctx context.Context, id string) ([]byte, error) {
	return nil, actorser.Ser().ActorMgr().KickLease(id)
}

func (d *ActorApp) SendKickLease(ctx context.Context, id string, app string) ([]byte, error) {
	return actorregcli.Kick(ctx, id, app)
}

func Run(ctx context.Context) error {
	loader.RegisterTitleCfg()
	runtime.New().Run(ctx, newApp(), []server.Server{actorser.New(), actorregser.New()})
	return nil
}
