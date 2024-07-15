package receiver

import (
	"github.com/sirupsen/logrus"
	"time"
)

// Receive инициализирует экземпляр RabbitMQ и постоянно прослушивает сообщения.
// В качестве входных данных для настройки соединения RabbitMQ принимается карта переменных среды.
// Функция устанавливает соединение с RabbitMQ, инициализирует очередь, а затем входит в бесконечный цикл для получения сообщений.
// При получении сообщения протоколируется тело сообщения.
// Соединение и канал закрываются с помощью defer, чтобы гарантировать их закрытие при выходе из функции или обнаружении ошибки.
func Receive(i interfaceReceiver, s interfaceStorage) {

	for {
		// Начинаем получать сообщения из очереди.
		msgs := i.ConsumeMessages()

		for d := range msgs {
			messageBody := d.Body
			messageID, errInsert := s.RecordMessage(messageBody)
			for errInsert != nil {
				logrus.Errorf("Failed to insertMessage: %s", errInsert)
				time.Sleep(5 * time.Second)
				messageID, errInsert = s.RecordMessage(messageBody)
			}
			logrus.Infof("Successful insert messageID = %v, messageBody = %s", messageID, messageBody)
		}
	}
}
