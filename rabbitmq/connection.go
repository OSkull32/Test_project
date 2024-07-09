package rabbitmq

import (
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
