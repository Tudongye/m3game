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
)

func CreateMapApp() *MapApp {
	return &MapApp{
		DefaultApp: app.CreateDefaultApp(dproto.MapAppFuncID),
	}
}

type MapApp struct {
	*app.DefaultApp
}

func (d *MapApp) Start() error {
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
	plugin.RegisterPluginFactory(&consul.Factory{})
	plugin.RegisterPluginFactory(&nats.Factory{})
	runtime.Run(CreateMapApp(), []server.Server{mapserver.CreateMapSer()})
	return nil
}
