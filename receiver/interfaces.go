package receiver

import amqp "github.com/rabbitmq/amqp091-go"

type interfaceReceiver interface {
	ConsumeMessages() <-chan amqp.Delivery
}

type interfaceStorage interface {
	RecordMessage(message []byte) (int, error)
}
