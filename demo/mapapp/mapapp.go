package mapapp

import (
	"m3game/app"
	"m3game/broker/nats"
	"m3game/demo/mapapp/mapserver"
	dproto "m3game/demo/proto"
	"m3game/mesh/router/consul"
	"m3game/runtime"
	"m3game/runtime/plugin"
	"m3game/server"
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
	router := plugin.GetRouterPlugin()
	if router != nil {
		if err := router.Register(d); err != nil {
			return err
		}
	}
	return nil
}
func (d *MapApp) HealthCheck() bool {
	return true
}
func Run() error {
	plugin.RegisterFactory(&consul.Factory{})
	plugin.RegisterFactory(&nats.Factory{})
	runtime.Run(newApp(), []server.Server{mapserver.New()})
	return nil
}
