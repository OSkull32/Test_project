package main

import (
	"os"
	"os/signal"
	"syscall"
	"test_project/config"

	"github.com/sirupsen/logrus"
	"test_project/receiver"
	"test_project/send"
)

func main() {
	log := logrus.New()

	env := config.LoadEnv()

	// Запускаем функцию Send() асинхронно
	go func() {
		log.Infoln("Starting Send()")
		send.Send(env) // Removed error handling
	}()

	// Запускаем функцию Receive() асинхронно
	go func() {
		log.Infoln("Starting Receive()")
		receiver.Receive(env) // Removed error handling
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, exiting...")
}
