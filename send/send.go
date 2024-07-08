package send

import (
	"context"
	"log"
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

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	rabbitmq.FailOnError(err, "Failed to declare a queue")

	for {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

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
		if err == nil {
			log.Printf(" [x] Sent %s\n", body)
		}
		rabbitmq.FailOnError(err, "Failed to publish a message")

		time.Sleep(10 * time.Second)
	}
}
