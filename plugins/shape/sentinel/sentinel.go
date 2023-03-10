package sentinel

import (
	"fmt"
	"m3game/plugins/log"
	"m3game/plugins/shape"
	"m3game/runtime/plugin"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	sentinelPlugin "github.com/sentinel-group/sentinel-go-adapters/grpc"
	"google.golang.org/grpc"
)

var (
	_         shape.Shape    = (*SentinelShape)(nil)
	_         plugin.Factory = (*Factory)(nil)
	_instance *SentinelShape
	_cfg      sentinelShapeCfg
	_factory  = &Factory{}
)

const (
	_factoryname = "shape_sentinel"
)

type sentinelShapeCfg struct {
	ConfigFile string `mapstructure:"ConfigFile"`
}

func init() {
	plugin.RegisterFactory(_factory)
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Shape
}

func (f *Factory) Name() string {
	return _factoryname
}

func (f *Factory) Setup(c map[string]interface{}) (plugin.PluginIns, error) {
	if _instance != nil {
		return _instance, nil
	}
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return nil, errors.Wrap(err, "Shape Decode Cfg")
	}
	if err := sentinel.InitWithConfigFile(_cfg.ConfigFile); err != nil {
		return nil, err
	}
	_instance = &SentinelShape{}
	shape.Set(_instance)
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

type SentinelShape struct {
}

func (s *SentinelShape) Factory() plugin.Factory {
	return _factory
}
func (s *SentinelShape) ClientInterceptor() grpc.UnaryClientInterceptor {
	return sentinelPlugin.NewUnaryClientInterceptor()
}
func (s *SentinelShape) ServerInterceptor() grpc.UnaryServerInterceptor {
	return sentinelPlugin.NewUnaryServerInterceptor()
}

func (s *SentinelShape) RegisterRule(rules []shape.Rule) error {
	var flowrules []*flow.Rule
	var breakules []*circuitbreaker.Rule
	for _, rule := range rules {
		log.Info("Register ShapeRule => %s", rule.Method)
		for _, frule := range rule.FlowRules {
			fr := &flow.Rule{
				Resource:         rule.Method,
				Threshold:        float64(frule.Threshold),
				StatIntervalInMs: uint32(frule.StatIntervalMs),
			}
			if frule.MaxQueueWaitMs != 0 {
				fr.MaxQueueingTimeMs = uint32(frule.MaxQueueWaitMs)
				fr.ControlBehavior = flow.Throttling
			}
			flowrules = append(flowrules, fr)
		}
		for _, brule := range rule.BreakRules {
			crule := &circuitbreaker.Rule{
				Resource:         rule.Method,
				RetryTimeoutMs:   uint32(brule.RetryTimeOutMs),
				StatIntervalMs:   uint32(brule.StatIntervalMs),
				MinRequestAmount: uint64(brule.MinRequestNum),
			}
			if shape.IsSlowRequst(brule.Strategy) {
				crule.Strategy = circuitbreaker.SlowRequestRatio
				crule.MaxAllowedRtMs = uint64(brule.SlowRequestMs)
				crule.Threshold = float64(brule.Threshold) / 1000
			} else if shape.IsErrorRatio(brule.Strategy) {
				crule.Strategy = circuitbreaker.ErrorRatio
				crule.Threshold = float64(brule.Threshold) / 1000
			} else if shape.IsErrorCount(brule.Strategy) {
				crule.Strategy = circuitbreaker.SlowRequestRatio
				crule.Threshold = float64(brule.Threshold)
			} else {
				return fmt.Errorf("Unknow Stategy %s", brule.Strategy)
			}
			breakules = append(breakules, crule)
		}
	}
	if _, err := flow.LoadRules(flowrules); err != nil {
		return err
	}
	if _, err := circuitbreaker.LoadRules(breakules); err != nil {
		return err
	}
	return nil
}
