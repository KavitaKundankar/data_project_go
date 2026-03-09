package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kavita/data-go-service/config"
	"github.com/kavita/data-go-service/internal/httpclient"
	"github.com/kavita/data-go-service/internal/processor"
	"github.com/kavita/data-go-service/internal/rabbitmq"
	"github.com/kavita/data-go-service/pkg/logger"
)

func main() {
	// Initialize Logger
	log := logger.NewLogger()
	log.Info("Starting Data Go Service...")

	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize RabbitMQ Connection
	rabbitConn, err := rabbitmq.NewConnection(cfg.RabbitMQURI)
	if err != nil {
		log.Error("Failed to initialize RabbitMQ: %v", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// Initialize HTTP Client
	httpClient := httpclient.NewClient()

	// Initialize Consumer
	consumer := rabbitmq.NewConsumer(rabbitConn, httpClient, log)

	// Start Consuming Queue 1 -> Function 1 -> HTTP 1
	err = consumer.StartConsuming(cfg.Queue1Name, cfg.HTTP1URL, processor.Function1)
	if err != nil {
		log.Error("Failed to start consumer for %s: %v", cfg.Queue1Name, err)
		os.Exit(1)
	}

	// Start Consuming Queue 2 -> Function 2 -> HTTP 2
	err = consumer.StartConsuming(cfg.Queue2Name, cfg.HTTP2URL, processor.Function2)
	if err != nil {
		log.Error("Failed to start consumer for %s: %v", cfg.Queue2Name, err)
		os.Exit(1)
	}

	log.Info("Service is running. Press CTRL+C to exit.")

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down service...")
}
