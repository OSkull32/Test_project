package send

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"test_project/rabbitmq"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Send() {
	conn, ch := rabbitmq.Connect()
	defer conn.Close()
	defer ch.Close()

	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		queueName = "hello"
	}

	for {

		q, err := ch.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		rabbitmq.FailOnError(err, "Send: Failed to declare a queue")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if conn.IsClosed() || ch.IsClosed() {
			logrus.Warn("Send: Connection or channel closed, attempting to reconnect...")
			conn, ch = rabbitmq.Connect() // Attempt to reconnect
			continue                      // Skip this iteration to retry in the next loop after reconnection
		}

		body := "Hello World!"
		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			logrus.Errorf("Failed to publish a message: %v", err)
			if conn.IsClosed() || ch.IsClosed() {
				logrus.Warn("Attempting to reconnect after failed publish...")
				conn, ch = rabbitmq.Connect() // Attempt to reconnect
			}
			continue // Skip this iteration to retry in the next loop after handling the error
		}

		logrus.Infof(" [x] Sent %s\n", body)
		time.Sleep(10 * time.Second)
	}
}
