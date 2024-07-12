package receiver

import (
	"github.com/sirupsen/logrus"
	"test_project/rabbitmq"
	"test_project/storage"
	"time"
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

	psqlDB := storage.InitPostgresDB(env)
	defer psqlDB.DB.Close()

	for {
		// Инициализируем очередь, чтобы начать получать сообщения.
		rabbitMQ.InitQueue()

		// Начинаем получать сообщения из очереди.
		msgs := rabbitMQ.ConsumeMessages()

		for d := range msgs {
			messageBody := d.Body
			messageID, errInsert := psqlDB.RecordMessage(messageBody)
			for errInsert != nil {
				logrus.Errorf("Failed to insertMessage: %s", errInsert)
				time.Sleep(5 * time.Second)
				messageID, errInsert = psqlDB.RecordMessage(messageBody)
			}
			logrus.Infof("Successful insert messageID = %v, messageBody = %s", messageID, messageBody)
		}
	}
}
