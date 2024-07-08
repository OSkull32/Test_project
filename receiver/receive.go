package receiver

import (
	"github.com/sirupsen/logrus"
	"os"
	"test_project/rabbitmq"
)

func Receive() {
	for {
		conn, ch := rabbitmq.Connect()

		defer func() {
			conn.Close()
			ch.Close()
		}()

		queueName := os.Getenv("QUEUE_NAME")
		if queueName == "" {
			queueName = "hello"
		}

		q, err := ch.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			logrus.Errorf("Failed to declare a queue: %v", err)
			continue
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			logrus.Errorf("Failed to register a consumer: %v", err)
			continue
		}

		logrus.Infof(" [*] Waiting for messages. To exit press CTRL+C")
		for d := range msgs {
			logrus.Infof("Received a message: %s", d.Body)
		}

	}
}
