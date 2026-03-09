package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func ConnectRabbitMQ(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Channel creation failed: %v", err)
	}

	return conn, ch
}
