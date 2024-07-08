package main

import (
	"fmt"
	"os"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func CheckErrSend(err error) {
	if err != nil {
		fmt.Printf("%s ", err)
		os.Exit(1)
	}
}
func main() {
	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost("localhost").
			SetPort(5552).
			SetUser("guest").
			SetPassword("guest"))
	CheckErrSend(err)

	streamName := "hello-go-stream"
	err = env.DeclareStream(streamName,
		&stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.GB(2),
		},
	)
	CheckErrSend(err)

	producer, err := env.NewProducer(streamName, stream.NewProducerOptions())
	CheckErrSend(err)

	err = producer.Send(amqp.NewMessage([]byte("Hello world")))
	CheckErrSend(err)
	fmt.Printf(" [x] 'Hello world' Message sent\n")
	err = producer.Close()
	CheckErrSend(err)
}
