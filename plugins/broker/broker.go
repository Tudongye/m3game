package broker

import (
	"m3game/meta/errs"
	"m3game/plugins/log"
	"m3game/runtime/plugin"
)

var (
	_broker Broker
)

type Broker interface {
	plugin.PluginIns
	Publish(topic string, bytes []byte) error
	Subscribe(topic string, handler func([]byte) error) error
}

func New(b Broker) (Broker, error) {
	if _broker != nil {
		log.Fatal("Broker Only One")
		return nil, errs.BrokerInsHasNewed.New("broker is newed %s", _broker.Factory().Name())
	}
	_broker = b
	return _broker, nil
}

func Instance() Broker {
	if _broker == nil {
		log.Fatal("Broker not newd")
		return nil
	}
	return _broker
}

func Publish(topic string, bytes []byte) error {
	if _broker == nil {
		return errs.BrokerInsIsNill.New("broker is nil")
	}
	return _broker.Publish(topic, bytes)
}

func Subscribe(topic string, h func([]byte) error) error {
	if _broker == nil {
		return errs.BrokerInsIsNill.New("broker is nil")
	}
	return _broker.Subscribe(topic, h)
}
