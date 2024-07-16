package main

import (
	"os"
	"os/signal"
	"syscall"
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

	r := rabbitmq.InitRabbitMQ(env)
	psqlDB := storage.InitPostgresDB(env)

	// Запускаем функцию отправки в новой горутине, чтобы запустить ее одновременно.
	go func() {
		log.Infoln("Starting Send()")
		send.Send(r)
	}()

	// Запускаем функцию приема в новой горутине, чтобы запустить ее одновременно.
	go func() {
		log.Infoln("Starting Receive()")
		receiver.Receive(r, psqlDB)
	}()

	// Обработка прерывания (Ctrl+C) и сигналов завершения для корректного завершения работы.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan // Block until a signal is received.
	log.Println("Shutdown signal received, exiting...")
}
