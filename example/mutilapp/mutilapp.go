package mutilapp

import (
	"m3game/example/mutilapp/mutilser"
	"m3game/example/proto"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
)

func newApp() *MutilApp {
	return &MutilApp{
		App: app.New(proto.MutilAppFuncID),
	}
}

type MutilApp struct {
	app.App
}

func (d *MutilApp) HealthCheck() bool {
	return true
}

func Run() error {
	runtime.Run(newApp(), []server.Server{mutilser.New()})
	return nil
}
