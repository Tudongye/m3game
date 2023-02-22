package dirapp

import (
	"m3game/app"
	"m3game/demo/dirapp/dirserver"
	"m3game/mesh/router/consul"
	"m3game/runtime"
	"m3game/runtime/plugin"
	"m3game/server"
)

const (
	AppFuncID = "dir"
)

func CreateDirApp() *DirApp {
	return &DirApp{
		DefaultApp: app.CreateDefaultApp(AppFuncID),
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
	runtime.Run(CreateDirApp(), []server.Server{dirserver.CreateDirSer()})
	return nil
}
