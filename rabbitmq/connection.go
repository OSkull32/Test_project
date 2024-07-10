package rabbitmq

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ(env map[string]string) *RabbitMQ {
	_, cancel := context.WithCancel(context.Background())
	res := &RabbitMQ{
		env:    env,
		cancel: cancel,
	}
	return res
}

type RabbitMQ struct {
	env      map[string]string
	cancel   func()
	ConnAmqp *amqp.Connection
	ChanAmqp *amqp.Channel
}

func FailOnError(err error, msg string) {
	if err != nil {
		logrus.Errorf("%s: %s", msg, err)
	}
}

func (r *RabbitMQ) Connect() {
	var err error
	r.ConnAmqp, err = r.tryConnect()
	for err != nil {
		logrus.Errorf("Failed to connect to RabbitMQ. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		r.ConnAmqp, err = r.tryConnect()
	}
	r.ChanAmqp, err = r.ConnAmqp.Channel()
	for err != nil {
		logrus.Errorf("Failed to open a channel. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		r.ChanAmqp, err = r.ConnAmqp.Channel()
	}
}

func (r *RabbitMQ) tryConnect() (*amqp.Connection, error) {
	UserName := r.env["RABBITMQ_DEFAULT_USER"]
	password := r.env["RABBITMQ_DEFAULT_PASS"]
	HOST := r.env["RABBITMQ_HOST"]
	PORT := r.env["RABBITMQ_PORT"]
	rabbitmqURL := &url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(UserName, password),
		Host:   fmt.Sprintf("%s:%s", HOST, PORT),
	}
	var err error
	r.ConnAmqp, err = amqp.Dial(rabbitmqURL.String())
	if err != nil {
		return nil, err
	}
	return r.ConnAmqp, nil
}

func (r *RabbitMQ) InitQueue() {
	queueName := r.env["QUEUE_NAME"]
	_, err := r.ChanAmqp.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	FailOnError(err, "Failed to declare a queue")
}

func (r *RabbitMQ) PublishMessage(ctx context.Context) {
	queueName := r.env["QUEUE_NAME"]
	body := r.env["MESSAGE_BODY"]
	err := r.ChanAmqp.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	FailOnError(err, "Failed to publish a message")
	logrus.Infof(" [x] Sent %s\n", body)
}

func (r *RabbitMQ) ConsumeMessages() (<-chan amqp.Delivery, error) {
	queueName := r.env["QUEUE_NAME"]
	msgs, err := r.ChanAmqp.Consume(
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
