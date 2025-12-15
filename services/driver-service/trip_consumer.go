package main

import (
	"context"
	"log"

	"ride-sharing/shared/messaging"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Listen() error
}

type tripConsumer struct {
	rabbitMQ *messaging.RabbitMQ
}

func NewTripConsumer(rabbitMQ *messaging.RabbitMQ) Consumer {
	return &tripConsumer{rabbitMQ}
}

func (c *tripConsumer) Listen() error {
	return c.rabbitMQ.ConsumeMessages(
		"hello",
		func(ctx context.Context, msg amqp.Delivery) error {
			log.Printf("Driver service received message: %v\n", msg)
			return nil
		},
	)
}
