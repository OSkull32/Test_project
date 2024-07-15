package rabbitmq

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/url"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// InitRabbitMQ инициализирует новую структуру RabbitMQ с предоставленными переменными среды.
// Он возвращает указатель на структуру RabbitMQ.
func InitRabbitMQ(env map[string]string) *RabbitMQ {
	ctx, cancel := context.WithCancel(context.Background())
	res := &RabbitMQ{
		env:    env,
		cancel: cancel,
		ctx:    ctx,
	}

	go res.ConnectRabbit(ctx)

	return res
}

// RabbitMQ содержит конфигурацию и состояние соединения RabbitMQ.
type RabbitMQ struct {
	ctx          context.Context
	env          map[string]string // Environment variables for RabbitMQ configuration.
	cancel       func()            // Function to cancel the context.
	ConnAmqp     *amqp.Connection  // Connection to RabbitMQ.
	PubChanAmqp  *amqp.Channel     // Channel for publishing messages
	ConsChanAmqp *amqp.Channel     // Channel for consuming messages
}

// FailOnError регистрирует сообщение об ошибке, если возникает ошибка.
func FailOnError(err error, msg string) {
	if err != nil {
		logrus.Errorf("%s: %s", msg, err)
	}
}

// ConnectRabbit пытается подключиться к RabbitMQ, используя переменные среды.
// Он повторяет попытку соединения в случае неудачи каждые 5 секунд.
func (r *RabbitMQ) ConnectRabbit(ctx context.Context) {
	for {
		if !r.CheckConnected() {
			var err error
			err = r.tryConnect(ctx)
			for err != nil {
				select {
				case <-ctx.Done():
					logrus.Errorf("Context cancelled, stopping connection attempts")
					return
				default:
					logrus.Errorf("Failed to connect to RabbitMQ. Retrying in 5 seconds...")
					time.Sleep(5 * time.Second)
					err = r.tryConnect(ctx)
				}
			}
			logrus.Info("Successful connection to RabbitMQ")
		}
	}
}

//
//// CreateChannel создает канал для общения с RabbitMQ.
//func (r *RabbitMQ) CreateChannel(ctx context.Context) {
//	var err error
//	if !r.CheckChannel() {
//		r.ChanAmqp, err = r.ConnAmqp.Channel()
//		for err != nil {
//			select {
//			case <-ctx.Done():
//				logrus.Errorf("Context cancelled, stopping channel opening attempts")
//				return
//			default:
//				logrus.Errorf("Failed to open a channel. Retrying in 5 seconds...")
//				time.Sleep(5 * time.Second)
//				r.ChanAmqp, err = r.ConnAmqp.Channel()
//			}
//		}
//		r.InitQueue()
//	}
//}

// CreatePubChannel creates a channel for publishing messages to RabbitMQ.
func (r *RabbitMQ) CreatePubChannel(ctx context.Context) {
	var err error
	if r.PubChanAmqp == nil || r.PubChanAmqp.IsClosed() {
		r.PubChanAmqp, err = r.ConnAmqp.Channel()
		for err != nil {
			select {
			case <-ctx.Done():
				logrus.Errorf("Context cancelled, stopping channel opening attempts")
				return
			default:
				logrus.Errorf("Failed to open a channel. Retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				r.PubChanAmqp, err = r.ConnAmqp.Channel()
			}
		}
		r.InitQueue()
	}
}

// CreateConsChannel creates a channel for consuming messages from RabbitMQ.
func (r *RabbitMQ) CreateConsChannel(ctx context.Context) {
	var err error
	if r.ConsChanAmqp == nil || r.ConsChanAmqp.IsClosed() {
		r.ConsChanAmqp, err = r.ConnAmqp.Channel()
		for err != nil {
			select {
			case <-ctx.Done():
				logrus.Errorf("Context cancelled, stopping channel opening attempts")
				return
			default:
				logrus.Errorf("Failed to open a channel. Retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				r.PubChanAmqp, err = r.ConnAmqp.Channel()
			}
		}
		r.InitQueue()
	}
}

// tryConnect пытается установить соединение с RabbitMQ, используя предоставленные учетные данные, и возвращает любую обнаруженную ошибку.
func (r *RabbitMQ) tryConnect(ctx context.Context) error {
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
	r.ConnAmqp, err = amqp.DialConfig(rabbitmqURL.String(), amqp.Config{
		Dial: func(network, addr string) (conn net.Conn, err error) {
			dialer := &net.Dialer{}
			return dialer.DialContext(ctx, network, addr) // Используем DialContext для учета контекста
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// InitQueue объявляет очередь на сервере RabbitMQ, используя переменные среды.
func (r *RabbitMQ) InitQueue() {
	queueName := r.env["QUEUE_NAME"]
	if r.ConsChanAmqp != nil && !r.ConsChanAmqp.IsClosed() {
		_, err := r.ConsChanAmqp.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		FailOnError(err, "Failed to declare a queue")
	} else {
		_, err := r.PubChanAmqp.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		FailOnError(err, "Failed to declare a queue")
	}
}

// PublishMessage публикует сообщение в очередь RabbitMQ, используя переменные среды.
// Для операции публикации используется контекст с тайм-аутом.
func (r *RabbitMQ) PublishMessage(ctx context.Context, body string) error {
	for !r.CheckConnected() {
		logrus.Warn("Connection closed, attempting to reconnect...")
		time.Sleep(5 * time.Second)
	}

	r.CreatePubChannel(r.ctx)

	queueName := r.env["QUEUE_NAME"]
	err := r.PubChanAmqp.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	return err
}

// ConsumeMessages начинает получать сообщения из указанной очереди.
// Он возвращает канал экземпляров доставки и все обнаруженные ошибки.
func (r *RabbitMQ) ConsumeMessages() <-chan amqp.Delivery {
	for !r.CheckConnected() {
		logrus.Warn("Connection closed, attempting to reconnect...")
		time.Sleep(5 * time.Second)
	}

	r.CreateConsChannel(r.ctx)
	queueName := r.env["QUEUE_NAME"]
	msg, err := r.ConsChanAmqp.Consume(
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
		return nil
	}
	return msg
}

func (r *RabbitMQ) CheckConnected() bool {
	return r.ConnAmqp != nil && !r.ConnAmqp.IsClosed()
}

//func (r *RabbitMQ) CheckChannel() bool {
//	return r.ChanAmqp != nil && !r.ChanAmqp.IsClosed()
//}

func (r *RabbitMQ) Shutdown() {
	if r != nil {
		r.cancel()
	}
}
