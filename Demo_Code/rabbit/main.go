package main

func main() {
	// Load configuration
	cfg := LoadConfig()

	// Connect to RabbitMQ
	conn, ch := ConnectRabbitMQ(cfg.RabbitMQURL)
	defer conn.Close()
	defer ch.Close()

	// Start consuming
	StartConsumer(ch, cfg.QueueName)
}
