package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"github.com/stretchr/testify/assert"
)

const (
	rabbitHost     = "localhost"
	rabbitPort     = 5552
	rabbitUser     = "guest"
	rabbitPassword = "guest"
	streamName     = "hello-go-stream"
	testMessage    = "Hello world"
)

func setupRabbitMQStream() (*stream.Environment, error) {
	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost(rabbitHost).
			SetPort(rabbitPort).
			SetUser(rabbitUser).
			SetPassword(rabbitPassword))
	if err != nil {
		return nil, err
	}

	err = env.DeclareStream(streamName,
		&stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.GB(2),
		},
	)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func TestSendMessageToStream(t *testing.T) {
	env, err := setupRabbitMQStream()
	assert.NoError(t, err)
	defer env.Close()

	producer, err := env.NewProducer(streamName, stream.NewProducerOptions())
	assert.NoError(t, err)
	defer producer.Close()

	err = producer.Send(amqp.NewMessage([]byte(testMessage)))
	assert.NoError(t, err)
	fmt.Printf(" [x] '%s' Message sent\n", testMessage)
}

func TestReceiveMessageFromStream(t *testing.T) {
	env, err := setupRabbitMQStream()
	assert.NoError(t, err)
	defer env.Close()

	messagesReceived := make(chan string)

	messagesHandler := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		messagesReceived <- string(byteSlicesToString(message.Data))
	}

	consumer, err := env.NewConsumer(streamName, messagesHandler,
		stream.NewConsumerOptions().SetOffset(stream.OffsetSpecification{}.First()))
	assert.NoError(t, err)
	defer consumer.Close()

	select {
	case msg := <-messagesReceived:
		assert.Equal(t, testMessage, msg)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for message")
	}

	fmt.Println(" [x] Message received successfully")
}

func byteSlicesToString(byteSlices [][]byte) string {
	str := make([]string, len(byteSlices))
	for i, b := range byteSlices {
		str[i] = string(b)
	}
	return strings.Join(str, "")
}
