package send

import (
	"context"
	"github.com/sirupsen/logrus"
	"test_project/rabbitmq"
	"time"
)

// Send инициализирует экземпляр RabbitMQ и непрерывно отправляет сообщения.
// В качестве входных данных для настройки соединения RabbitMQ принимается карта переменных среды.
// Функция устанавливает соединение с RabbitMQ, инициализирует очередь, а затем входит в цикл для отправки сообщений.
// Если соединение или канал закрыты, он пытается повторно подключиться, прежде чем продолжить отправку сообщений.
func Send(env map[string]string) {
	// Создание экземпляра RabbitMQ
	rabbitMQ := rabbitmq.InitRabbitMQ(env)

	// Соединение и канал закрыты при выходе из функции или возникновении ошибки.
	defer func() {
		rabbitMQ.ConnAmqp.Close()
		rabbitMQ.ChanAmqp.Close()
	}()

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
