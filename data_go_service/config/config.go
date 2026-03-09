package config

import (
	"log"
	"os"
)

type Config struct {
	RabbitMQURI string
	Queue1Name  string
	Queue2Name  string
	HTTP1URL    string
	HTTP2URL    string
}

func LoadConfig() *Config {
	// In a real app, you might use a library like viper or godotenv
	// For simplicity, we'll use os.Getenv with defaults

	rabbitURI := getEnv("RABBITMQ_URL", "amqp://user:password@localhost:5672/myvhost")
	q1 := getEnv("queue1", "dummy_data_queue_1")
	q2 := getEnv("queue2", "dummy_data_queue_2")
	h1 := getEnv("http1", "http://localhost:8080/api1")
	h2 := getEnv("http2", "http://localhost:8080/api2")

	return &Config{
		RabbitMQURI: rabbitURI,
		Queue1Name:  q1,
		Queue2Name:  q2,
		HTTP1URL:    h1,
		HTTP2URL:    h2,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Warning: Environment variable %s not set, using default: %s", key, fallback)
	return fallback
}
