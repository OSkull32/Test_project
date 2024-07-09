package send

import (
	"context"
	"github.com/sirupsen/logrus"
	"test_project/rabbitmq"
	"time"
)

func Send(env map[string]string) {
	conn, ch := rabbitmq.Connect(env)
	defer conn.Close()
	defer ch.Close()

	queueName := env["QUEUE_NAME"]

	for {

		rabbitmq.InitQueue(ch, queueName)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if conn.IsClosed() || ch.IsClosed() {
			logrus.Warn("Send: Connection or channel closed, attempting to reconnect...")
			conn, ch = rabbitmq.Connect(env) // Attempt to reconnect
			continue                         // Skip this iteration to retry in the next loop after reconnection
		}
		body := "Hello World!"
		rabbitmq.PublishMessage(ch, queueName, body, ctx)

		logrus.Infof(" [x] Sent %s\n", body)
		time.Sleep(10 * time.Second)
	}
}
