// Package events provides event publishing and consuming functionalities
package events

import (
	"context"
	"time"

	"ride-sharing/shared/messaging"
)

type TripEventsPublisher struct {
	rabbitMQ *messaging.RabbitMQ
}

func NewPublisher(rabbitMQ *messaging.RabbitMQ) Publisher {
	return &TripEventsPublisher{
		rabbitMQ: rabbitMQ,
	}
}

func (t *TripEventsPublisher) Publish(ctx context.Context, event string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return t.rabbitMQ.PublishWithContext(ctx, "update_me", event)
}
