package prometheus

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/meta/errs"
	"m3game/plugins/log"
	"m3game/plugins/metric"
	"m3game/plugins/router"
	"m3game/runtime/plugin"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	_        metric.Metric    = (*Metric)(nil)
	_        plugin.Factory   = (*Factory)(nil)
	_        plugin.PluginIns = (*Metric)(nil)
	_metric  *Metric
	_factory = &Factory{}
)

const (
	_name = "metric_prom"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type PromCfg struct {
	Host            string `mapstructure:"Host" validate:"required"` // 监听地址
	Port            int    `mapstructure:"Port" validate:"gt=0"`     // 监听端口
	ConsulSvc       string `mapstructure:"ConsulSvc"`
	ConsulAppPrefix string `mapstructure:"ConsulAppPrefix"`
}
type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Metric
}
func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(ctx context.Context, c map[string]interface{}) (plugin.PluginIns, error) {
	if _metric != nil {
		return _metric, nil
	}
	var cfg PromCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errs.PromSetupFaul.Wrap(err, "Router Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, errs.PromSetupFaul.Wrap(err, "")
	}
	_metric = &Metric{}
	if _, err := metric.New(_metric); err != nil {
		return nil, errs.PromSetupFaul.Wrap(err, "")
	}

	listenport := fmt.Sprintf(":%d", cfg.Port)
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Metric.Listen %s...", listenport)
	go func() {
		if err := http.ListenAndServe(listenport, nil); err != nil {
			log.Fatal(err.Error())
		}
	}()

	if cfg.ConsulSvc != "" {
		if err := router.Instance().Register(
			fmt.Sprintf("%s.%s", cfg.ConsulAppPrefix, config.GetAppID().String()),
			cfg.ConsulSvc,
			cfg.Host,
			cfg.Port,
			nil,
			func(string, string) bool {
				return true
			},
		); err != nil {
			return nil, errs.PromSetupFaul.Wrap(err, "")
		}
	}

	return _metric, nil
}

func (f *Factory) Destroy(plugin.PluginIns) error {
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanUnload(plugin.PluginIns) bool {
	return false
}

type Metric struct {
}

func (*Metric) Factory() plugin.Factory {
	return _factory
}

func (*Metric) NewCounter(key string, group string) metric.StatCounter {
	return promauto.NewCounter(
		prometheus.CounterOpts{
			Subsystem: group,
			Name:      key,
		})
}
func (*Metric) NewGauge(key string, group string) metric.StatGauge {
	return promauto.NewGauge(
		prometheus.GaugeOpts{
			Subsystem: group,
			Name:      key,
		})
}

func (*Metric) NewHistogram(key string, group string) metric.StatHistogram {
	return promauto.NewHistogram(
		prometheus.HistogramOpts{
			Subsystem: group,
			Name:      key,
		})
}

func (*Metric) NewSummary(key string, group string) metric.StatSummary {
	return promauto.NewSummary(
		prometheus.SummaryOpts{
			Subsystem:  group,
			Name:       key,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})
}
