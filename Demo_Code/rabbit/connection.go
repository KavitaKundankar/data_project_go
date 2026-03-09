package main

import (
	"log"

	"github.com/streadway/amqp"
)

func ConnectRabbitMQ(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("RabbitMQ connection failed:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Channel creation failed:", err)
	}

	return conn, ch
}
