package mapapp

import (
	"m3game/app"
	"m3game/demo/mapapp/mapserver"
	"m3game/mesh/router/consul"
	"m3game/runtime"
	"m3game/runtime/plugin"
	"m3game/server"
)

const (
	AppFuncID = "map"
)

func CreateMapApp() *MapApp {
	return &MapApp{
		DefaultApp: app.CreateDefaultApp(AppFuncID),
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
	runtime.Run(CreateMapApp(), []server.Server{mapserver.CreateMapSer()})
	return nil
}
