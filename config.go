package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func LoadEnv() map[string]string {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	requiredEnv := []string{"RABBITMQ_DEFAULT_USER", "RABBITMQ_DEFAULT_PASS", "RABBITMQ_HOST", "RABBITMQ_PORT", "QUEUE_NAME"}
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
