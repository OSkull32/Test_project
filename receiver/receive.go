package receiver

import (
	"github.com/sirupsen/logrus"
	"test_project/rabbitmq"
)

func Receive(env map[string]string) {
	for {
		conn, ch := rabbitmq.Connect(env)
		defer func() {
			conn.Close()
			ch.Close()
		}()

		queueName := env["QUEUE_NAME"]
		rabbitmq.InitQueue(ch, queueName)

		msgs, err := rabbitmq.ConsumeMessages(ch, queueName)
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
