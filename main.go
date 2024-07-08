package main

import (
	"github.com/joho/godotenv"
	"log"
	"test_project/receiver"
	"test_project/send"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Запускаем функцию Send() асинхронно
	go func() {
		log.Println("Starting Send()")
		send.Send()
	}()

	// Запускаем функцию Receive() асинхронно
	go func() {
		log.Println("Starting Receive()")
		receiver.Receive()
	}()
	select {}
}
