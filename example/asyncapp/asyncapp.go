package asyncapp

import (
	"m3game/example/asyncapp/asyncser"
	"m3game/example/proto"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
)

func newApp() *AsyncApp {
	return &AsyncApp{
		App: app.New(proto.AsyncAppFuncID),
	}
}

type AsyncApp struct {
	app.App
}

func (d *AsyncApp) HealthCheck() bool {
	return true
}

func Run() error {
	runtime.Run(newApp(), []server.Server{asyncser.New()})
	return nil
}
