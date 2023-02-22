package roleapp

import (
	"m3game/app"
	"m3game/db/cache"
	"m3game/demo/dirapp/dirclient"
	"m3game/demo/loader"
	"m3game/demo/mapapp/mapclient"
	"m3game/demo/roleapp/roleserver"
	"m3game/mesh/router/consul"
	"m3game/proto"
	"m3game/runtime"
	"m3game/runtime/plugin"
	"m3game/server"
)

const (
	AppFuncID = "role"
)

func CreateRoleApp() *RoleApp {
	return &RoleApp{
		DefaultApp: app.CreateDefaultApp(AppFuncID),
	}
}

type RoleApp struct {
	*app.DefaultApp
}

func (d *RoleApp) Start() error {
	router := plugin.GetRouterPlugin()
	if router != nil {
		if err := router.Register(d); err != nil {
			return err
		}
	}
	if err := dirclient.Init(d.RouteIns(), func(c *dirclient.Client) { c.Client = proto.META_FLAG_FALSE }); err != nil {
		return err
	}
	if err := mapclient.Init(d.RouteIns(), func(c *mapclient.Client) { c.Client = proto.META_FLAG_FALSE }); err != nil {
		return err
	}
	return nil
}
func (d *RoleApp) HealthCheck() bool {
	return true
}
func Run() error {
	plugin.RegisterPluginFactory(&consul.Factory{})
	plugin.RegisterPluginFactory(&cache.Factory{})

	loader.RegisterLocationCfg()

	runtime.Run(CreateRoleApp(), []server.Server{roleserver.CreateRoleSer()})
	return nil
}
