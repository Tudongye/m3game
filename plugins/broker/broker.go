package broker

import (
	"fmt"
	"m3game/runtime/plugin"
)

var (
	_broker Broker
)

type Broker interface {
	plugin.PluginIns
	Publish(topic string, m []byte) error
	Subscribe(topic string, h func([]byte)) error
}

func GenTopic(c string) string {
	return fmt.Sprintf("Topic.%s", c)
}
func Set(r Broker) {
	if _broker != nil {
		panic("Broker Only One")
	}
	_broker = r
}

func Get() Broker {
	if _broker == nil {
		panic("Broker Mush Have One")
	}
	return _broker
}

func Publish(topic string, m []byte) error {
	return Get().Publish(topic, m)
}

func Subscribe(topic string, h func([]byte)) error {
	return Get().Subscribe(topic, h)
}
