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
)

func CreateDirApp() *DirApp {
	return &DirApp{
		DefaultApp: app.CreateDefaultApp(dproto.DirAppFuncID),
	}
}

type DirApp struct {
	*app.DefaultApp
}

func (d *DirApp) Start() error {
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
	plugin.RegisterPluginFactory(&consul.Factory{})
	plugin.RegisterPluginFactory(&nats.Factory{})
	runtime.Run(CreateDirApp(), []server.Server{dirserver.CreateDirSer()})
	return nil
}
