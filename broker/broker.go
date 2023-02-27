package broker

type Broker interface {
	Publish(topic string, m []byte) error
	Subscribe(topic string, h func([]byte)) error
}
