package consul

import (
	"fmt"
	"log"
	"m3game/mesh/router"
	"m3game/runtime/plugin"
	"m3game/runtime/transport"
	"m3game/util"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/mapstructure"
)

var (
	_         router.Router  = (*Router)(nil)
	_         plugin.Factory = (*Factory)(nil)
	_cfg                     = consulRouterCfg{}
	_instance *Router
)

const (
	_factoryname = "router_consul"
)

type consulRouterCfg struct {
	ConsulHost string
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Router
}
func (f *Factory) Name() string {
	return _factoryname
}

func (f *Factory) Setup(c map[string]interface{}) (plugin.PluginIns, error) {
	if _instance != nil {
		return _instance, nil
	}
	_instance = &Router{
		apps: make(map[string]router.AppReciver),
	}
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return nil, err
	}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = _cfg.ConsulHost
	if client, err := api.NewClient(consulConfig); err != nil {
		return nil, err
	} else {
		_instance.client = client
	}
	return _instance, nil
}

func (f *Factory) Destroy(plugin.PluginIns) error {
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanDelete(plugin.PluginIns) bool {
	return false
}

type Router struct {
	client *api.Client
	apps   map[string]router.AppReciver
}

func (r *Router) Name() string {
	return _factoryname
}

func (r *Router) Register(app router.AppReciver) error {
	r.apps[app.IDStr()] = app
	envid, worldid, funcid, _, err := util.AppStr2ID(app.IDStr())
	if err != nil {
		return err
	}
	svcstr := util.SvcID2Str(envid, worldid, funcid)
	ip, port, err := util.Addr2IPPort(transport.Addr())
	if err != nil {
		return err
	}
	interval := time.Duration(10) * time.Second
	deregister := time.Duration(1) * time.Minute
	reg := &api.AgentServiceRegistration{
		ID:      app.IDStr(), // 服务节点的名称
		Name:    svcstr,      // 服务名称
		Tags:    []string{},  // tag，可以为空
		Port:    port,        // 服务端口
		Address: ip,          // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval:                       interval.String(),                                     // 健康检查间隔
			GRPC:                           fmt.Sprintf("%v:%v/Health/%v", ip, port, app.IDStr()), // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: deregister.String(),                                   // 注销时间，相当于过期时间
		},
	}
	log.Println(fmt.Sprintf("%v:%v/Health/%v", ip, port, app.IDStr()))
	agent := r.client.Agent()
	if err := agent.ServiceRegister(reg); err != nil {
		return err
	}
	return nil
}

func (r *Router) Deregister(app router.AppReciver) error {
	delete(r.apps, app.IDStr())
	return nil
}

func (r *Router) GetAllInstances(envid string, worldid string, funcid string) ([]router.Instance, error) {
	return nil, nil
}
