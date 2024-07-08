package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"test_project/receiver"
	"test_project/send"
)

func main() {
	log := logrus.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Запускаем функцию Send() асинхронно
	go func() {
		log.Infoln("Starting Send()")
		send.Send() // Removed error handling
	}()

	// Запускаем функцию Receive() асинхронно
	go func() {
		log.Infoln("Starting Receive()")
		receiver.Receive() // Removed error handling
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, exiting...")
}
