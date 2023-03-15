package actorapp

import (
	"context"
	"m3game/config"
	"m3game/example/actorapp/actor"
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
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func newApp() *ActorApp {
	return &ActorApp{
		App:  app.New(proto.ActorAppFuncID),
		exit: make(chan struct{}, 1),
	}
}

var (
	_cfg AppCfg
)

type ActorApp struct {
	app.App
	exit chan struct{}
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

func (a *ActorApp) Init(cfg map[string]interface{}) error {
	if err := mapstructure.Decode(cfg, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return err
	}
	return nil
}
func (d *ActorApp) Start(wg *sync.WaitGroup) error {
	if err := asynccli.Init(config.GetAppID()); err != nil {
		return err
	}
	if err := gatecli.Init(config.GetAppID()); err != nil {
		return err
	}
	if err := actorregcli.Init(config.GetAppID()); err != nil {
		return err
	}
	lease.SetReciver(d)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("ActorApp PrepareTime %d", _cfg.PrePareTime)
		time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
		log.Info("ActorApp Ready")
		t := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-d.exit:
				return
			case <-t.C:
				// 插件检查
				if router.Get().Factory().CanDelete(router.Get()) {
					t.Stop()
					runtime.ShutDown()
				}
				if lease.Get().Factory().CanDelete(lease.Get()) {
					t.Stop()
					runtime.ShutDown()
				}
				continue
			}
		}
	}()
	return nil
}

func (d *ActorApp) Stop() error {
	select {
	case d.exit <- struct{}{}:
		return nil
	default:
		return nil
	}
}

func (d *ActorApp) HealthCheck() bool {
	return true
}

func (d *ActorApp) PreExitLease(ctx context.Context, id string) ([]byte, error) {
	actorid := actor.ParseActorIdFromLeaseId(id)
	return nil, actorser.Ser().ActorMgr().KickOne(actorid)
}

func (d *ActorApp) SendKickLease(ctx context.Context, id string, app string) ([]byte, error) {
	return actorregcli.Kick(ctx, id, app)
}

func Run() error {
	loader.RegisterTitleCfg()
	runtime.Run(newApp(), []server.Server{actorser.New(), actorregser.New()})
	return nil
}
