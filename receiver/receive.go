package receiver

import (
	"github.com/sirupsen/logrus"
	"test_project/rabbitmq"
)

func Receive(env map[string]string) {
	for {
		// Создание экземпляра RabbitMQ
		rabbitMQ := rabbitmq.InitRabbitMQ(env)
		// Вызов метода Connect
		conn, ch := rabbitMQ.Connect()
		defer func() {
			conn.Close()
			ch.Close()
		}()

		rabbitMQ.InitQueue()

		msgs, err := rabbitMQ.ConsumeMessages()
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
