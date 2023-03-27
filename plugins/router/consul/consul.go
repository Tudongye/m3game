package consul

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/plugins/log"
	"m3game/plugins/router"
	"m3game/runtime/plugin"
	"m3game/util"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	_             router.Router  = (*Router)(nil)
	_             plugin.Factory = (*Factory)(nil)
	_consulrouter *Router
	_factory      = &Factory{}
)

func init() {
	plugin.RegisterFactory(_factory)
}

const (
	_name = "router_consul"
)

type consulRouterCfg struct {
	ConsulHost string `mapstructure:"ConsulHost" validate:"required"`
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Router
}
func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(ctx context.Context, c map[string]interface{}) (plugin.PluginIns, error) {
	if _consulrouter != nil {
		return _consulrouter, nil
	}
	var cfg consulRouterCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errors.Wrap(err, "Router Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, err
	}
	_consulrouter = &Router{}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.ConsulHost
	if client, err := api.NewClient(consulConfig); err != nil {
		return nil, errors.Wrapf(err, "Consul.Api.NewClient Add %s", consulConfig.Address)
	} else {
		_consulrouter.client = client
	}

	if _, err := router.New(_consulrouter); err != nil {
		return nil, err
	}
	return _consulrouter, nil
}

func (f *Factory) Destroy(p plugin.PluginIns) error {
	r := p.(*Router)
	r.Deregister(config.GetAppID().String(), config.GetSvcID().String())
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanUnload(p plugin.PluginIns) bool {
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

func (r *Router) Register(app string, svc string, addr string, meta map[string]string) error {
	ip, port, err := util.Addr2IPPort(addr)
	if err != nil {
		return err
	}
	interval := time.Duration(10) * time.Second
	deregister := time.Duration(1) * time.Minute
	healthmethod := fmt.Sprintf("%v:%v/Health/%v", ip, port, app)
	if meta == nil {
		meta = make(map[string]string)
	}
	reg := &api.AgentServiceRegistration{
		ID:      app,        // 服务节点的名称
		Name:    svc,        // 服务名称
		Tags:    []string{}, // tag，可以为空
		Port:    port,       // 服务端口
		Address: ip,         // 服务 IP
		Meta:    meta,
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

func (r *Router) GetAllInstances(svcid string) ([]router.Ins, error) {
	services, _, err := _consulrouter.client.Health().Service(svcid, "", true, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Consul.Health.Service %s fail", svcid)
	}
	var instances []router.Ins
	for _, service := range services {
		instances = append(instances, newInstance(service.Service.Address, service.Service.Port, service.Service.ID, service.Service.Meta))
	}
	return instances, nil
}
