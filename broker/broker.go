package broker

import "fmt"

type Broker interface {
	Publish(topic string, m []byte) error
	Subscribe(topic string, h func([]byte)) error
}

func GenTopic(c string) string {
	return fmt.Sprintf("Topic.%s", c)
}
