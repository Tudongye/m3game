package multiapp

import (
	"m3game/example/multiapp/multiser"
	"m3game/example/proto"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
)

func newApp() *MultiApp {
	return &MultiApp{
		App: app.New(proto.MultiAppFuncID),
	}
}

type MultiApp struct {
	app.App
}

func (d *MultiApp) HealthCheck() bool {
	return true
}

func Run() error {
	runtime.Run(newApp(), []server.Server{multiser.New()})
	return nil
}
