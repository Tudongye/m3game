package prometheus

import (
	"fmt"
	"m3game/config"
	"m3game/plugins/log"
	"m3game/plugins/metric"
	"m3game/runtime/plugin"
	"m3game/runtime/transport"
	"m3game/util"
	"net/http"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	_         metric.Metric    = (*Metric)(nil)
	_         plugin.Factory   = (*Factory)(nil)
	_         plugin.PluginIns = (*Metric)(nil)
	_instance *Metric
	_cfg      PromCfg
	_factory  = &Factory{}
)

const (
	_factoryname = "metric_prom"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Metric
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

	_instance = &Metric{}
	metric.Set(_instance)
	listenport := fmt.Sprintf(":%d", _cfg.Port)
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Metric.Listen %s...", listenport)
	go func() {
		if err := http.ListenAndServe(listenport, nil); err != nil {
			panic(err.Error())
		}
	}()

	if _cfg.ConsulUrl != "" {
		if err := registerConsul(_cfg.ConsulUrl, config.GetSvcID().String(), config.GetAppID().String(), _cfg.Port); err != nil {
			return nil, err
		}
		log.Info("Metric.registerConsul %s svc %s ins %s...", _cfg.ConsulUrl, config.GetSvcID(), config.GetAppID())
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

type PromCfg struct {
	Port      int    `mapstructure:"Port"`       // Prom监听端口
	ConsulUrl string `mapstructure:"ConsulHost"` // Consul注册
}

func (p *PromCfg) CheckVaild() error {
	if p.Port == 0 {
		return errors.New("PromCfg.Port cannot = 0")
	}
	return nil
}

type Metric struct {
}

func (*Metric) Factory() plugin.Factory {
	return _factory
}

func (*Metric) NewCounter(key string, group string) metric.StatCounter {
	return &pstatCounter{
		c: promauto.NewCounter(
			prometheus.CounterOpts{
				Subsystem: group,
				Name:      key,
			}),
	}
}
func (*Metric) NewGauge(key string, group string) metric.StatGauge {
	return &pstatGauge{
		c: promauto.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: group,
				Name:      key,
			}),
	}
}

func (*Metric) NewHistogram(key string, group string) metric.StatHistogram {
	return &pstatHistogram{
		c: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: group,
				Name:      key,
			}),
	}
}

func (*Metric) NewSummary(key string, group string) metric.StatSummary {
	return &pstatSummary{
		c: promauto.NewSummary(
			prometheus.SummaryOpts{
				Subsystem:  group,
				Name:       key,
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			}),
	}
}

func registerConsul(consulurl string, svc string, ins string, port int) error {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulurl
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return errors.Wrap(err, "Metric.registerConsul")
	}
	interval := time.Duration(10) * time.Second
	deregister := time.Duration(1) * time.Minute
	addr := transport.Addr()
	ip, _, err := util.Addr2IPPort(addr)
	if err != nil {
		return err
	}
	reg := &api.AgentServiceRegistration{
		ID:      ins,        // 服务节点的名称
		Name:    svc,        // 服务名称
		Tags:    []string{}, // tag，可以为空
		Port:    port,       // 服务端口
		Address: ip,         // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval:                       interval.String(), // 健康检查间隔
			TCP:                            addr,
			DeregisterCriticalServiceAfter: deregister.String(), // 注销时间，相当于过期时间
		},
	}
	agent := client.Agent()
	if err := agent.ServiceRegister(reg); err != nil {
		return errors.Wrapf(err, "Metric.registerConsul")
	}
	return nil
}
