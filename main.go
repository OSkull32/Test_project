package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"test_project/receiver"
	"test_project/send"
)

func main() {
	log := logrus.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Запускаем функцию Send() асинхронно
	go func() {
		log.Infoln("Starting Send()")
		send.Send()
	}()

	// Запускаем функцию Receive() асинхронно
	go func() {
		log.Infoln("Starting Receive()")
		receiver.Receive()
	}()
	select {}
}
