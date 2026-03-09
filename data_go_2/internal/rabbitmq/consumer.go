package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func Consume(queueName string, handler func([]byte), ch *amqp.Channel) {
	msgs, err := ch.Consume(
		queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer for %s: %v", queueName, err)
	}

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()
}
