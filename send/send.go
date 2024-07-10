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
	rabbitMQ.Connect()
	defer rabbitMQ.ConnAmqp.Close()
	defer rabbitMQ.ChanAmqp.Close()

	for {
		rabbitMQ.InitQueue()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if rabbitMQ.ConnAmqp.IsClosed() || rabbitMQ.ChanAmqp.IsClosed() {
			logrus.Warn("Send: Connection or channel closed, attempting to reconnect...")
			rabbitMQ.Connect() // Используем метод Connect экземпляра
			continue
		}

		rabbitMQ.PublishMessage(ctx)

		time.Sleep(10 * time.Second)
	}
}
