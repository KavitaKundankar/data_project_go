package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"data_go_3/internal/models"
	"data_go_3/internal/processor"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	processor *processor.Processor
}

func NewConsumer(url string, p *processor.Processor) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &Consumer{
		conn:      conn,
		channel:   ch,
		processor: p,
	}, nil
}

func (c *Consumer) Start(ctx context.Context, queue1Name, queue2Name string) error {
	q1, err := c.channel.QueueDeclare(queue1Name, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue 1: %w", err)
	}

	q2, err := c.channel.QueueDeclare(queue2Name, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue 2: %w", err)
	}

	msgs1, err := c.channel.Consume(q1.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register consumer for queue 1: %w", err)
	}

	msgs2, err := c.channel.Consume(q2.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register consumer for queue 2: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for d := range msgs1 {
			c.processMessage(d, 1)
		}
	}()

	go func() {
		defer wg.Done()
		for d := range msgs2 {
			c.processMessage(d, 2)
		}
	}()

	log.Println("RabbitMQ Consumers started...")

	<-ctx.Done()
	return nil
}

func (c *Consumer) processMessage(d amqp.Delivery, queueIndex int) {
	var body map[string]interface{}
	if err := json.Unmarshal(d.Body, &body); err != nil {
		log.Printf("Error unmarshaling message from Queue %d: %v", queueIndex, err)
		d.Nack(false, false)
		return
	}

	// We pass the entire body as the payload to keep it simple and as received
	data := models.QueueData{
		ID:      d.MessageId,
		Payload: body,
	}

	var processed map[string]interface{}
	var err error

	if queueIndex == 1 {
		processed, err = c.processor.ProcessQueue1(data)
		if err == nil {
			err = c.processor.SendToQueue1API(processed)
		}
	} else {
		processed, err = c.processor.ProcessQueue2(data)
		if err == nil {
			err = c.processor.SendToQueue2API(processed)
		}
	}

	if err != nil {
		log.Printf("Error processing message from Queue %d: %v", queueIndex, err)
		d.Nack(false, true) // Requeue on error
	} else {
		d.Ack(false)
		log.Printf("Message from Queue %d processed and sent.", queueIndex)
	}
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
