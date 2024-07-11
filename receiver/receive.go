package receiver

import (
	"github.com/sirupsen/logrus"
	"test_project/rabbitmq"
	"test_project/storage"
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

	psqlDB := storage.NewPsqlDB(env)

	defer psqlDB.Close()

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

		for d := range msgs {
			message, errInsert := storage.InsertMessage(psqlDB, env["MESSAGE_BODY"])
			if errInsert != nil {
				logrus.Errorf("InsertMessage: %s", err)
			}
			if errInsert == nil {
				logrus.Infof("Successful insert messageID = %v, messageBody = %s", message, d.Body)
			}
		}
	}
}
