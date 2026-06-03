package broker

type MessageBroker interface {
	Publish(queue string, body []byte) error
	Subscribe(queue string, handler func([]byte) error) error
}
