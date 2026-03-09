package config

import (
	"os"
	"strconv"
)

type Config struct {
	RabbitMQURL  string
	Queue1Name   string
	Queue2Name   string
	HTTP1URL     string
	HTTP2URL     string
	APIPort      int
}

func LoadConfig() *Config {
	return &Config{
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://user:password@localhost:5672/myvhost"),
		Queue1Name:  getEnv("QUEUE_1_NAME", "dummy_data_queue_1"),
		Queue2Name:  getEnv("QUEUE_2_NAME", "dummy_data_queue_2"),
		HTTP1URL:    getEnv("HTTP_1_URL", "http://localhost:8000/api/data1"),
		HTTP2URL:    getEnv("HTTP_2_URL", "http://localhost:8000/api/data2"),
		APIPort:     getEnvInt("API_PORT", 8000),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
