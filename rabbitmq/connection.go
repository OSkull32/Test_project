package rabbitmq

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		logrus.Panicf("%s: %s", msg, err)
	}
}

func Connect() (*amqp.Connection, *amqp.Channel) {
	conn, err := tryConnect()
	for err != nil {
		logrus.Errorf("Failed to connect to RabbitMQ. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		conn, err = tryConnect()
	}
	ch, err := conn.Channel()
	for err != nil {
		logrus.Errorf("Failed to open a channel. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		ch, err = conn.Channel()
	}
	return conn, ch
}

func tryConnect() (*amqp.Connection, error) {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
