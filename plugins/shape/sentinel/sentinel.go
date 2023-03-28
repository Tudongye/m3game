package sentinel

import (
	"context"
	"m3game/meta/errs"
	"m3game/plugins/log"
	"m3game/plugins/shape"
	"m3game/runtime/plugin"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	sentinelPlugin "github.com/sentinel-group/sentinel-go-adapters/grpc"
	"google.golang.org/grpc"
)

var (
	_              shape.Shape    = (*SentinelShape)(nil)
	_              plugin.Factory = (*Factory)(nil)
	_sentinelshape *SentinelShape
	_factory       = &Factory{}
)

const (
	_name = "shape_sentinel"
)

type sentinelShapeCfg struct {
	ConfigFile string `mapstructure:"ConfigFile" validate:"required"`
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
	return _name
}

func (f *Factory) Setup(ctx context.Context, c map[string]interface{}) (plugin.PluginIns, error) {
	if _sentinelshape != nil {
		return _sentinelshape, nil
	}
	var cfg sentinelShapeCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errs.SentinelSteupFail.Wrap(err, "Shape Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, errs.SentinelSteupFail.Wrap(err, "")
	}
	if err := sentinel.InitWithConfigFile(cfg.ConfigFile); err != nil {
		return nil, errs.SentinelSteupFail.Wrap(err, "")
	}
	_sentinelshape = &SentinelShape{}
	if _, err := shape.New(_sentinelshape); err != nil {
		return nil, err
	}
	return _sentinelshape, nil
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
			flowRule := &flow.Rule{
				Resource:         rule.Method,
				Threshold:        float64(frule.Threshold),
				StatIntervalInMs: uint32(frule.StatIntervalMs),
			}
			if frule.MaxQueueWaitMs != 0 {
				flowRule.MaxQueueingTimeMs = uint32(frule.MaxQueueWaitMs)
				flowRule.ControlBehavior = flow.Throttling
			}
			flowrules = append(flowrules, flowRule)
		}
		for _, brule := range rule.BreakRules {
			circuitBreakerRule := &circuitbreaker.Rule{
				Resource:         rule.Method,
				RetryTimeoutMs:   uint32(brule.RetryTimeOutMs),
				StatIntervalMs:   uint32(brule.StatIntervalMs),
				MinRequestAmount: uint64(brule.MinRequestNum),
			}
			if shape.IsSlowRequst(brule.Strategy) {
				circuitBreakerRule.Strategy = circuitbreaker.SlowRequestRatio
				circuitBreakerRule.MaxAllowedRtMs = uint64(brule.SlowRequestMs)
				circuitBreakerRule.Threshold = float64(brule.Threshold) / 1000
			} else if shape.IsErrorRatio(brule.Strategy) {
				circuitBreakerRule.Strategy = circuitbreaker.ErrorRatio
				circuitBreakerRule.Threshold = float64(brule.Threshold) / 1000
			} else if shape.IsErrorCount(brule.Strategy) {
				circuitBreakerRule.Strategy = circuitbreaker.SlowRequestRatio
				circuitBreakerRule.Threshold = float64(brule.Threshold)
			} else {
				return errs.SentinelRegisterRuleFail.New("Unknow Stategy %s", brule.Strategy)
			}
			breakules = append(breakules, circuitBreakerRule)
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
