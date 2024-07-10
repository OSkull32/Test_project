package send

import (
	"context"
	"github.com/sirupsen/logrus"
	"test_project/rabbitmq"
	"time"
)

func Send(env map[string]string) {
	// Создание экземпляра RabbitMQ
	rabbitMQ := rabbitmq.InitRabbitMQ(env)

	// Вызов метода Connect
	conn, ch := rabbitMQ.Connect()
	defer conn.Close()
	defer ch.Close()

	for {
		rabbitMQ.InitQueue()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if conn.IsClosed() || ch.IsClosed() {
			logrus.Warn("Send: Connection or channel closed, attempting to reconnect...")
			conn, ch = rabbitMQ.Connect() // Используем метод Connect экземпляра
			continue
		}

		rabbitMQ.PublishMessage(ctx)

		time.Sleep(10 * time.Second)
	}
}
