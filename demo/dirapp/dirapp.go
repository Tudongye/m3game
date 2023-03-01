package dirapp

import (
	"m3game/app"
	"m3game/broker/nats"
	"m3game/demo/dirapp/dirserver"
	dproto "m3game/demo/proto"
	"m3game/mesh/router/consul"
	"m3game/runtime"
	"m3game/runtime/plugin"
	"m3game/server"
	"sync"
)

func newApp() *DirApp {
	return &DirApp{
		App: app.New(dproto.DirAppFuncID),
	}
}

type DirApp struct {
	app.App
}

func (d *DirApp) Start(wg *sync.WaitGroup) error {
	router := plugin.GetRouterPlugin()
	if router != nil {
		if err := router.Register(d); err != nil {
			return err
		}
	}
	return nil
}
func (d *DirApp) HealthCheck() bool {
	return true
}
func Run() error {
	plugin.RegisterFactory(&consul.Factory{})
	plugin.RegisterFactory(&nats.Factory{})
	runtime.Run(newApp(), []server.Server{dirserver.New()})
	return nil
}
