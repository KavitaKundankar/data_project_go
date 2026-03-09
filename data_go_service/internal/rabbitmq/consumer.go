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
func (c *Consumer) StartConsuming(queueName string, httpURL string, processFunc func(map[string]interface{}) models.ProcessedResult) error {
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
			c.log.Info("-------------------------------------------")
			c.log.Info("EVENT: Received a message from %s", queueName)

			var payload map[string]interface{}
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				c.log.Error("ERROR: Failed to unmarshal JSON: %v", err)
				d.Nack(false, false)
				continue
			}

			c.log.Info("DEBUG: Message Content: %v", payload)

			// 1. Process data
			result := processFunc(payload)

			c.log.Info("ACTION: Sending processed data to %s", httpURL)

			// 2. Send via HTTP
			err = c.httpClient.SendData(httpURL, result)
			if err != nil {
				c.log.Error("ERROR: HTTP request failed: %v", err)
				c.log.Info("REQUEUING: Requeuing for %s", queueName)
				d.Nack(false, true)
				continue
			}

			// 3. Acknowledge message only after HTTP success
			c.log.Info("DONE: Acknowledging message from %s", queueName)
			d.Ack(false)
			c.log.Info("-------------------------------------------")
		}
	}()

	return nil
}
