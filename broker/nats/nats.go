package nats

import (
	"m3game/broker"
	"m3game/log"
	"m3game/runtime/plugin"

	"github.com/mitchellh/mapstructure"
	"github.com/nats-io/nats.go"
)

var (
	_         broker.Broker    = (*Broker)(nil)
	_         plugin.PluginIns = (*Broker)(nil)
	_         plugin.Factory   = (*Factory)(nil)
	_cfg                       = natsBrokerCfg{}
	_instance *Broker
)

const (
	_factoryname = "broker_nats"
)

func init() {
	plugin.RegisterFactory(&Factory{})
}

type natsBrokerCfg struct {
	NatsURL string
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.Broker
}
func (f *Factory) Name() string {
	return _factoryname
}

func (f *Factory) Setup(c map[string]interface{}) (plugin.PluginIns, error) {
	if _instance != nil {
		return _instance, nil
	}
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return nil, err
	}
	_instance = &Broker{
		subs: make(map[string]*nats.Subscription),
	}
	if nc, err := nats.Connect(_cfg.NatsURL); err != nil {
		return nil, err
	} else {
		_instance.nc = nc
		if js, err := nc.JetStream(nats.PublishAsyncMaxPending(256)); err != nil {
			return nil, err
		} else {
			_instance.js = js
		}
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

type Broker struct {
	nc   *nats.Conn
	js   nats.JetStreamContext
	subs map[string]*nats.Subscription
}

func (b *Broker) Publish(topic string, m []byte) error {
	_, err := b.js.PublishAsync(topic, m)
	return err
}

func (b *Broker) Subscribe(topic string, h func([]byte)) error {
	log.Info("Subscribe %s", topic)
	_, err := b.nc.Subscribe(topic, func(m *nats.Msg) {
		h(m.Data)
	})
	return err
}

func (b *Broker) Name() string {
	return _factoryname
}
