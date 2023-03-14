package actorapp

import (
	"m3game/config"
	"m3game/example/actorapp/actorregser"
	"m3game/example/actorapp/actorser"
	"m3game/example/asyncapp/asynccli"
	"m3game/example/gateapp/gatecli"
	"m3game/example/loader"
	"m3game/example/proto"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/db/cache"
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"sync"
)

func newApp() *ActorApp {
	return &ActorApp{
		App: app.New(proto.ActorAppFuncID),
	}
}

type ActorApp struct {
	app.App
}

func (d *ActorApp) Start(wg *sync.WaitGroup) error {
	if err := asynccli.Init(config.GetAppID()); err != nil {
		return err
	}
	if err := gatecli.Init(config.GetAppID()); err != nil {
		return err
	}

	return nil
}
func (d *ActorApp) HealthCheck() bool {
	return true
}
func Run() error {
	loader.RegisterTitleCfg()
	runtime.Run(newApp(), []server.Server{actorser.New(), actorregser.New()})
	return nil
}
