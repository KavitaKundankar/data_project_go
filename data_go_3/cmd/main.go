package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"data_go_3/internal/api"
	"data_go_3/internal/config"
	"data_go_3/internal/processor"
	"data_go_3/internal/rabbitmq"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize processor
	proc := processor.NewProcessor(cfg.HTTP1URL, cfg.HTTP2URL)

	// Initialize API server
	srv := api.NewServer(cfg.APIPort)

	// Initialize RabbitMQ consumer
	cons, err := rabbitmq.NewConsumer(cfg.RabbitMQURL, proc)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}
	defer cons.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start RabbitMQ consumers in background
	go func() {
		if err := cons.Start(ctx, cfg.Queue1Name, cfg.Queue2Name); err != nil {
			log.Printf("Consumer stopped with error: %v", err)
			cancel()
		}
	}()

	// Start API server (this is blocking)
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("API server stopped with error: %v", err)
			cancel()
		}
	}()

	log.Println("Microservice is running. Press Ctrl+C to exit.")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Received signal %v. Shutting down...", sig)
	case <-ctx.Done():
		log.Println("Context cancelled. Shutting down...")
	}
}
