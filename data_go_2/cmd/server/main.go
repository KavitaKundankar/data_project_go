package main

import (
	"go-rabbit-microservice/config"
	"go-rabbit-microservice/internal/processor"
	"go-rabbit-microservice/internal/rabbitmq"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to RabbitMQ
	conn, ch := rabbitmq.ConnectRabbitMQ(cfg.RabbitMQURL)
	defer conn.Close()
	defer ch.Close()

	log.Println("Microservice Started")

	// Initialize processors with target HTTP URLs
	p1 := processor.NewProcessor(cfg.HTTPService1)

	// Start consuming from queues
	rabbitmq.Consume(cfg.Queue1, p1.Process, ch)

	// For queue 2, we can use a wrapper to pass the URL
	rabbitmq.Consume(cfg.Queue2, func(data []byte) {
		processor.ProcessQueue2(cfg.HTTPService2, data)
	}, ch)

	log.Println("Consuming from queues:", cfg.Queue1, "and", cfg.Queue2)

	// Block forever
	select {}
}
