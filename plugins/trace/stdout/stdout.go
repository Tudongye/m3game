package stdout

import (
	"m3game/plugins/trace"
	"m3game/runtime/plugin"

	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	_         plugin.Factory = (*Factory)(nil)
	_instance *Trace
	_factory  = &Factory{}
)

const (
	_name = "trace_stdout"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Trace
}

func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(c map[string]interface{}) (plugin.PluginIns, error) {
	if _instance != nil {
		return _instance, nil
	}
	_instance = &Trace{}
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
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

type Trace struct {
}

func (s *Trace) Factory() plugin.Factory {
	return _factory
}
