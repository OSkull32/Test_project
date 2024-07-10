package rabbitmq

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// InitRabbitMQ инициализирует новую структуру RabbitMQ с предоставленными переменными среды.
// Он возвращает указатель на структуру RabbitMQ.
func InitRabbitMQ(env map[string]string) *RabbitMQ {
	_, cancel := context.WithCancel(context.Background())
	res := &RabbitMQ{
		env:    env,
		cancel: cancel,
	}
	return res
}

// RabbitMQ содержит конфигурацию и состояние соединения RabbitMQ.
type RabbitMQ struct {
	env      map[string]string // Environment variables for RabbitMQ configuration.
	cancel   func()            // Function to cancel the context.
	ConnAmqp *amqp.Connection  // Connection to RabbitMQ.
	ChanAmqp *amqp.Channel     // Channel for RabbitMQ communication.
}

// FailOnError регистрирует сообщение об ошибке, если возникает ошибка.
func FailOnError(err error, msg string) {
	if err != nil {
		logrus.Errorf("%s: %s", msg, err)
	}
}

// Connect пытается подключиться к RabbitMQ, используя переменные среды.
// Он повторяет попытку соединения в случае неудачи каждые 5 секунд.
func (r *RabbitMQ) Connect() {
	var err error
	err = r.tryConnect()
	for err != nil {
		logrus.Errorf("Failed to connect to RabbitMQ. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		err = r.tryConnect()
	}
	r.ChanAmqp, err = r.ConnAmqp.Channel()
	for err != nil {
		logrus.Errorf("Failed to open a channel. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		r.ChanAmqp, err = r.ConnAmqp.Channel()
	}
}

// tryConnect пытается установить соединение с RabbitMQ, используя предоставленные учетные данные, и возвращает любую обнаруженную ошибку.
func (r *RabbitMQ) tryConnect() error {
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
		return err
	}
	return nil
}

// InitQueue объявляет очередь на сервере RabbitMQ, используя переменные среды.
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

// PublishMessage публикует сообщение в очередь RabbitMQ, используя переменные среды.
// Для операции публикации используется контекст с тайм-аутом.
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

// ConsumeMessages начинает получать сообщения из указанной очереди.
// Он возвращает канал экземпляров доставки и все обнаруженные ошибки.
func (r *RabbitMQ) ConsumeMessages() (<-chan amqp.Delivery, error) {
	queueName := r.env["QUEUE_NAME"]
	msg, err := r.ChanAmqp.Consume(
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
	return msg, nil
}
