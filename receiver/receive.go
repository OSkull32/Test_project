package receiver

import (
	"github.com/sirupsen/logrus"
	"test_project/rabbitmq"
)

// Receive инициализирует экземпляр RabbitMQ и постоянно прослушивает сообщения.
// В качестве входных данных для настройки соединения RabbitMQ принимается карта переменных среды.
// Функция устанавливает соединение с RabbitMQ, инициализирует очередь, а затем входит в бесконечный цикл для получения сообщений.
// При получении сообщения протоколируется тело сообщения.
// Соединение и канал закрываются с помощью defer, чтобы гарантировать их закрытие при выходе из функции или обнаружении ошибки.
func Receive(env map[string]string) {
	// Создание экземпляра RabbitMQ
	rabbitMQ := rabbitmq.InitRabbitMQ(env)

	// Соединение и канал закрыты при выходе из функции или возникновении ошибки.
	defer func() {
		rabbitMQ.ConnAmqp.Close()
		rabbitMQ.ChanAmqp.Close()
	}()

	for {
		// Инициализируем очередь, чтобы начать получать сообщения.
		rabbitMQ.InitQueue()

		if rabbitMQ.ConnAmqp.IsClosed() || rabbitMQ.ChanAmqp.IsClosed() {
			logrus.Warn("Send: Connection or channel closed, attempting to reconnect...")
			rabbitMQ.Connect(rabbitMQ.Ctx) // Используем метод Connect экземпляра
			continue
		}

		// Начинаем получать сообщения из очереди.
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
