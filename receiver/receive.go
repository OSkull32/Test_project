package receiver

import (
	"github.com/sirupsen/logrus"
	"os"
	"test_project/rabbitmq"
)

func Receive() {
	conn, ch := rabbitmq.Connect()
	defer conn.Close()
	defer ch.Close()

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
	rabbitmq.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	rabbitmq.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			logrus.Infof("Received a message: %s", d.Body)
		}
	}()

	logrus.Infof(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
