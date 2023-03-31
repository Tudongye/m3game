package jaeger

import (
	"context"
	"m3game/config"
	"m3game/meta/errs"
	"m3game/plugins/trace"
	"m3game/runtime/plugin"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	_         plugin.Factory = (*Factory)(nil)
	_instance *JaegerTracer
	_factory  = &Factory{}
)

const (
	_name = "trace_jaeger"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type JaegerCfg struct {
	Host string `mapstructure:"Host" validate:"required"`
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Trace
}

func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(ctx context.Context, c map[string]interface{}) (plugin.PluginIns, error) {
	if _instance != nil {
		return _instance, nil
	}
	var cfg JaegerCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errs.MongoSetupFail.Wrap(err, "JaegerCfg Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, errs.MongoSetupFail.Wrap(err, "")
	}
	_instance = &JaegerTracer{}
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(cfg.Host)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.GetAppID().String()),
			attribute.String("Env", config.GetEnvID().String()),
			attribute.String("World", config.GetWorldID().String()),
			attribute.String("Svc", config.GetSvcID().String()),
			attribute.String("App", config.GetAppID().String()),
		),
		),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	if _, err := trace.New(tp); err != nil {
		return nil, err
	}
	return _instance, nil
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

type JaegerTracer struct {
}

func (s *JaegerTracer) Factory() plugin.Factory {
	return _factory
}
