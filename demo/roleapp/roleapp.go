package roleapp

import (
	"m3game/app"
	"m3game/broker/nats"
	"m3game/db/cache"
	"m3game/demo/dirapp/dirclient"
	"m3game/demo/loader"
	"m3game/demo/mapapp/mapclient"
	dproto "m3game/demo/proto"
	"m3game/demo/roleapp/rolechclient"
	"m3game/demo/roleapp/rolechserver"
	"m3game/demo/roleapp/roleserver"
	"m3game/mesh/router/consul"
	"m3game/proto"
	"m3game/runtime"
	"m3game/runtime/plugin"
	"m3game/server"
)

func CreateRoleApp() *RoleApp {
	return &RoleApp{
		DefaultApp: app.CreateDefaultApp(dproto.RoleAppFuncID),
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
	if err := rolechclient.Init(d.RouteIns(), func(c *rolechclient.Client) { c.Client = proto.META_FLAG_FALSE }); err != nil {
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
	plugin.RegisterPluginFactory(&nats.Factory{})

	loader.RegisterLocationCfg()

	runtime.Run(CreateRoleApp(), []server.Server{roleserver.CreateRoleSer(), rolechserver.CreateRoleChSer()})
	return nil
}
