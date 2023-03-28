package nats

import (
	"context"
	"m3game/meta/errs"
	"m3game/plugins/broker"
	"m3game/plugins/log"
	"m3game/runtime/plugin"

	"github.com/go-playground/validator/v10"

	"github.com/mitchellh/mapstructure"
	"github.com/nats-io/nats.go"
)

var (
	_        broker.Broker    = (*Broker)(nil)
	_        plugin.PluginIns = (*Broker)(nil)
	_        plugin.Factory   = (*Factory)(nil)
	_broker  *Broker
	_factory = &Factory{}
)

const (
	_name = "broker_nats"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type natsBrokerCfg struct {
	URL string `mapstructure:"URL" validate:"required"`
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Broker
}
func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(ctx context.Context, c map[string]interface{}) (plugin.PluginIns, error) {
	if _broker != nil {
		return _broker, nil
	}
	var cfg natsBrokerCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errs.NatsSetupFail.Wrap(err, "")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, errs.NatsSetupFail.Wrap(err, "")
	}
	_broker = &Broker{
		subs: make(map[string]*nats.Subscription),
	}
	if nc, err := nats.Connect(cfg.URL); err != nil {
		return nil, errs.NatsSetupFail.Wrap(err, "Nats.Conntect %s", cfg.URL)
	} else {
		_broker.nc = nc
		if js, err := nc.JetStream(nats.PublishAsyncMaxPending(256)); err != nil {
			return nil, errs.NatsSetupFail.Wrap(err, "nc.JetStream %s", cfg.URL)
		} else {
			_broker.js = js
		}
	}
	if _, err := broker.New(_broker); err != nil {
		return nil, err
	}
	return _broker, nil
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

type Broker struct {
	nc   *nats.Conn
	js   nats.JetStreamContext
	subs map[string]*nats.Subscription
}

func (b *Broker) Factory() plugin.Factory {
	return _factory
}

func (b *Broker) Publish(topic string, m []byte) error {
	_, err := b.js.PublishAsync(topic, m)
	return err
}

func (b *Broker) Subscribe(topic string, h func([]byte) error) error {
	log.Info("Subscribe %s", topic)
	_, err := b.nc.Subscribe(topic, func(m *nats.Msg) {
		if err := h(m.Data); err != nil {
			log.Error("broker subscribe %s handler err %s", topic, err.Error())
		}
	})
	return err
}
