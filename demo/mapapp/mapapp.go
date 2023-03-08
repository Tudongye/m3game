package mapapp

import (
	_ "m3game/broker/nats"
	"m3game/demo/mapapp/mapserver"
	dproto "m3game/demo/proto"
	"m3game/mesh/router"
	_ "m3game/mesh/router/consul"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	_ "m3game/shape/sentinel"
	"sync"
)

func newApp() *MapApp {
	return &MapApp{
		App: app.New(dproto.MapAppFuncID),
	}
}

type MapApp struct {
	app.App
}

func (d *MapApp) Start(wg *sync.WaitGroup) error {
	if err := router.Register(d); err != nil {
		return err
	}
	return nil
}
func (d *MapApp) HealthCheck() bool {
	return true
}
func Run() error {
	runtime.Run(newApp(), []server.Server{mapserver.New()})
	return nil
}
