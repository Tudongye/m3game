package consul

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/meta/errs"
	"m3game/plugins/log"
	"m3game/plugins/router"
	"m3game/runtime/plugin"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/mapstructure"
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
	AliveHost  string `mapstructure:"AliveHost" validate:"required"` // 本地心跳检测IP
	AlivePort  int    `mapstructure:"AlivePort" validate:"gt=0"`     // 本地心跳检测Port
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
		return nil, errs.ConsulSetupFail.Wrap(err, "Router Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, errs.ConsulSetupFail.Wrap(err, "")
	}
	_consulrouter = &Router{
		cfg: cfg,
	}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.ConsulHost
	if client, err := api.NewClient(consulConfig); err != nil {
		return nil, errs.ConsulSetupFail.Wrap(err, "Consul.Api.NewClient Add %s", consulConfig.Address)
	} else {
		_consulrouter.client = client
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		app := queryValues.Get("app")
		svc := queryValues.Get("svc")
		if !_consulrouter.aliveCheck(app, svc) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	server := &http.Server{Addr: ":" + strconv.Itoa(cfg.AlivePort)}
	go func() {
		log.Debug("Consul Alive HTTP server listening on %v", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Consul Alive HTTP server failed: %v", err)
		}
	}()
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
			if ins.GetAppID() == string(config.GetAppID()) {
				return false
			}
		}
		return true
	}
}

type Router struct {
	client *api.Client
	cfg    consulRouterCfg
	livefs sync.Map
}

func (r *Router) Factory() plugin.Factory {
	return _factory
}

func (r *Router) Register(app string, svc string, host string, port int, meta map[string]string, livef func(app string, svc string) bool) error {
	interval := time.Duration(10) * time.Second
	deregister := time.Duration(1) * time.Minute
	healthmethod := fmt.Sprintf("http://%v:%v/health?svc=%v&app=%v", r.cfg.AliveHost, r.cfg.AlivePort, svc, app)
	if meta == nil {
		meta = make(map[string]string)
	}
	reg := &api.AgentServiceRegistration{
		ID:      app,        // 服务节点的名称
		Name:    svc,        // 服务名称
		Tags:    []string{}, // tag，可以为空
		Port:    port,       // 服务端口
		Address: host,       // 服务 IP
		Meta:    meta,
		Check: &api.AgentServiceCheck{ // 健康检查
			HTTP:                           healthmethod,
			Method:                         "GET",
			Interval:                       interval.String(),
			Timeout:                        "3s",
			DeregisterCriticalServiceAfter: deregister.String(),
		},
	}
	log.Info("Register HealthMethod => %s", healthmethod)
	agent := r.client.Agent()
	if err := agent.ServiceRegister(reg); err != nil {
		return errs.ConsulRegisterAppFail.Wrap(err, "Consul.agent.ServiceRegister %s", app)
	}
	r.livefs.Store(app, livef)

	return nil
}

func (r *Router) Deregister(app string, svc string) error {
	return nil
}

func (r *Router) GetAllInstances(svcid string) ([]router.Ins, error) {
	services, _, err := _consulrouter.client.Health().Service(svcid, "", true, nil)
	if err != nil {
		return nil, errs.ConsulGetAllInstanceFail.Wrap(err, "Consul.Health.Service %s fail", svcid)
	}
	var instances []router.Ins
	for _, service := range services {
		instances = append(instances, newInstance(service.Service.Address, service.Service.Port, service.Service.ID, service.Service.Meta))
	}
	return instances, nil
}

func (r *Router) aliveCheck(app string, svc string) bool {
	if v, ok := r.livefs.Load(app); !ok {
		return false
	} else if f, ok := v.(func(app string, svc string) bool); !ok {
		return false
	} else {
		return f(app, svc)
	}
}
