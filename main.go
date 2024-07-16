package main

import (
	"context"
	"test_project/config"
	"test_project/rabbitmq"
	"test_project/storage"

	"github.com/sirupsen/logrus"
	"test_project/receiver"
	"test_project/send"
)

// main — это точка входа приложения.
// Он инициализирует настройки среды, запускает функции отправки и получения в отдельных горутинах,
// и ждет сигнала прерывания или завершения, чтобы корректно завершить работу приложения.
func main() {
	log := logrus.New() // Инициализируем новый logger.

	env := config.LoadEnv() // Загрузка переменных среды.
	globalCtx, globalCancel := context.WithCancel(context.Background())
	rabbitCtx, rabbit := rabbitmq.InitRabbitMQ(globalCtx, env)
	psqlCtx, psqlDB := storage.InitPostgresDB(globalCtx, env)

	defer func() {
		log.Infoln("Shutting down...")
		rabbit.Shutdown()
		psqlDB.Shutdown()
		globalCancel()
	}()

	// Запускаем функцию отправки в новой горутине, чтобы запустить ее одновременно.
	go func() {
		log.Infoln("Starting Send()")
		send.Send(rabbit)
	}()

	// Запускаем функцию приема в новой горутине, чтобы запустить ее одновременно.
	go func() {
		log.Infoln("Starting Receive()")
		receiver.Receive(rabbit, psqlDB)
	}()

	for {
		select {
		case <-globalCtx.Done():
			log.Infoln("globalCtx is done")
			return
		case <-rabbitCtx.Done():
			log.Infoln("rabbitCtx is done")
			return
		case <-psqlCtx.Done():
			log.Infoln("psqlCtx is done")
			return
		}

	}
}
