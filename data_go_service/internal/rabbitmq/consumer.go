package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/kavita/data-go-service/internal/httpclient"
	"github.com/kavita/data-go-service/internal/models"
	"github.com/kavita/data-go-service/pkg/logger"
)

type Consumer struct {
	rabbitConn *Connection
	httpClient *httpclient.Client
	log        *logger.Logger
}

func NewConsumer(conn *Connection, client *httpclient.Client, log *logger.Logger) *Consumer {
	return &Consumer{
		rabbitConn: conn,
		httpClient: client,
		log:        log,
	}
}

// StartConsuming starts taking messages from a queue and applies processing and http sending
func (c *Consumer) StartConsuming(queueName string, httpURL string, processFunc func(models.DataPayload) models.ProcessedResult) error {
	msgs, err := c.rabbitConn.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (IMPORTANT: false so we ack manually)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			c.log.Info("Received a message from %s", queueName)

			var payload models.DataPayload
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				c.log.Error("Error unmarshaling JSON: %v", err)
				// You might want to Nack here if the message is invalid
				d.Nack(false, false)
				continue
			}

			// 1. Process data (Event based)
			result := processFunc(payload)

			// 2. Send via HTTP
			err = c.httpClient.SendData(httpURL, result)
			if err != nil {
				c.log.Error("Error sending HTTP request to %s: %v", httpURL, err)

				// [FIX]: Check if it's a connection error or a transient error
				// For this demo, we'll requeue after a short delay or just nack and discard if it keeps failing
				// To keep it simple: Let's log it clearly and requeue, but alert the user.
				c.log.Info("Requeuing message for %s. Suggestion: Start the mock server!", queueName)
				d.Nack(false, true)
				continue
			}

			// 3. Acknowledge message only after HTTP success
			d.Ack(false)
			c.log.Info("Successfully processed and acknowledged message from %s", queueName)
		}
	}()

	return nil
}
