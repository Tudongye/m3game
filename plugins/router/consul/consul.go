package consul

import (
	"fmt"
	"m3game/config"
	"m3game/plugins/log"
	"m3game/plugins/router"
	"m3game/runtime/plugin"
	"m3game/util"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	_         router.Router  = (*Router)(nil)
	_         plugin.Factory = (*Factory)(nil)
	_cfg                     = consulRouterCfg{}
	_instance *Router
	_factory  = &Factory{}
)

func init() {
	plugin.RegisterFactory(_factory)
}

const (
	_factoryname = "router_consul"
)

type consulRouterCfg struct {
	ConsulHost string `mapstructure:"ConsulHost"`
}

func (c *consulRouterCfg) CheckVaild() error {
	if c.ConsulHost == "" {
		return errors.New("ConsulHost cant be space")
	}
	return nil
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
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return nil, errors.Wrap(err, "Router Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return nil, err
	}
	_instance = &Router{}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = _cfg.ConsulHost
	if client, err := api.NewClient(consulConfig); err != nil {
		return nil, errors.Wrapf(err, "Consul.Api.NewClient Add %s", consulConfig.Address)
	} else {
		_instance.client = client
	}

	router.Set(_instance)
	return _instance, nil
}

func (f *Factory) Destroy(p plugin.PluginIns) error {
	r := p.(*Router)
	r.Deregister(config.GetAppID().String(), config.GetSvcID().String())
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanDelete(p plugin.PluginIns) bool {
	r := p.(*Router)
	if inss, err := r.GetAllInstances(config.GetSvcID().String()); err != nil {
		return true
	} else if len(inss) == 0 {
		return true
	} else {
		for _, ins := range inss {
			if ins.GetIDStr() == string(config.GetAppID()) {
				return false
			}
		}
		return true
	}
}

type Router struct {
	client *api.Client
}

func (r *Router) Factory() plugin.Factory {
	return _factory
}

func (r *Router) Register(app string, svc string, addr string) error {
	ip, port, err := util.Addr2IPPort(addr)
	if err != nil {
		return err
	}
	interval := time.Duration(10) * time.Second
	deregister := time.Duration(1) * time.Minute
	healthmethod := fmt.Sprintf("%v:%v/Health/%v", ip, port, app)
	reg := &api.AgentServiceRegistration{
		ID:      app,        // 服务节点的名称
		Name:    svc,        // 服务名称
		Tags:    []string{}, // tag，可以为空
		Port:    port,       // 服务端口
		Address: ip,         // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval:                       interval.String(),   // 健康检查间隔
			GRPC:                           healthmethod,        // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: deregister.String(), // 注销时间，相当于过期时间
		},
	}
	log.Info("Register HealthMethod => %s", healthmethod)
	agent := r.client.Agent()
	if err := agent.ServiceRegister(reg); err != nil {
		return errors.Wrapf(err, "Consul.agent.ServiceRegister %s", app)
	}
	return nil
}

func (r *Router) Deregister(app string, svc string) error {
	return nil
}

func (r *Router) GetAllInstances(svcid string) ([]router.Instance, error) {
	services, _, err := _instance.client.Health().Service(svcid, "", true, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Consul.Health.Service %s fail", svcid)
	}
	var instances []router.Instance
	for _, service := range services {
		instances = append(instances, newInstance(service.Service.Address, service.Service.Port, service.Service.ID))
	}
	return instances, nil
}
