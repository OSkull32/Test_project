package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

// LoadEnv загружает переменные среды из файла .env и обеспечивает установку всех необходимых переменных среды.
// Он возвращает карту необходимых переменных среды и их значений.
// Если файл .env не может быть загружен или какая-либо необходимая переменная среды отсутствует, приложение завершится.
//
// Обязательные переменные среды:
// - RABBITMQ_DEFAULT_USER: имя пользователя для аутентификации RabbitMQ.
// - RABBITMQ_DEFAULT_PASS: пароль для аутентификации RabbitMQ.
// - RABBITMQ_HOST: имя хоста сервера RabbitMQ.
// - RABBITMQ_PORT: порт, на котором работает сервер RabbitMQ.
// - QUEUE_NAME: имя используемой очереди RabbitMQ.
// - MESSAGE_BODY: тело сообщения
func LoadEnv() map[string]string {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	requiredEnv := []string{"RABBITMQ_DEFAULT_USER", "RABBITMQ_DEFAULT_PASS", "RABBITMQ_HOST", "RABBITMQ_PORT", "QUEUE_NAME", "MESSAGE_BODY"}
	envMap := make(map[string]string)
	for _, envKey := range requiredEnv {
		value := os.Getenv(envKey)
		if value == "" {
			logrus.Fatalf("Environment variable %s is required", envKey)
		}
		envMap[envKey] = value
	}
	return envMap
}
