package send

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

// Send инициализирует экземпляр RabbitMQ и непрерывно отправляет сообщения.
// В качестве входных данных для настройки соединения RabbitMQ принимается карта переменных среды.
// Функция устанавливает соединение с RabbitMQ, инициализирует очередь, а затем входит в цикл для отправки сообщений.
// Если соединение или канал закрыты, он пытается повторно подключиться, прежде чем продолжить отправку сообщений.
func Send(i interfaceSend) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		message := "Hello World!"

		err := i.PublishMessage(ctx, message)
		if err != nil {
			logrus.Errorf("Failed to insertMessage: %s", err)
			time.Sleep(5 * time.Second)
		} else {
			logrus.Infof("Successful send message: %s", message)
		}

		time.Sleep(10 * time.Second)
	}
}
