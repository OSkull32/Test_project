package rabbitmq

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		logrus.Errorf("%s: %s", msg, err)
	}
}

func Connect(env map[string]string) (*amqp.Connection, *amqp.Channel) {
	conn, err := tryConnect(env)
	for err != nil {
		logrus.Errorf("Failed to connect to RabbitMQ. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		conn, err = tryConnect(env)
	}
	ch, err := conn.Channel()
	for err != nil {
		logrus.Errorf("Failed to open a channel. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		ch, err = conn.Channel()
	}
	return conn, ch
}

func tryConnect(env map[string]string) (*amqp.Connection, error) {
	UserName := env["RABBITMQ_DEFAULT_USER"]
	password := env["RABBITMQ_DEFAULT_PASS"]
	HOST := env["RABBITMQ_HOST"]
	PORT := env["RABBITMQ_PORT"]
	rabbitmqURL := &url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(UserName, password),
		Host:   fmt.Sprintf("%s:%s", HOST, PORT),
	}

	conn, err := amqp.Dial(rabbitmqURL.String())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func InitQueue(ch *amqp.Channel, queueName string) {
	_, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	FailOnError(err, "Failed to declare a queue")
}

func PublishMessage(ch *amqp.Channel, queueName string, body string, ctx context.Context) {
	err := ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		logrus.Errorf("Failed to publish a message: %v", err)
	}
}

func ConsumeMessages(ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		logrus.Errorf("Failed to register a consumer: %v", err)
		return nil, err
	}
	return msgs, nil
}
