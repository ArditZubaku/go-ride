// Package messaging
package messaging

import (
	"ride-sharing/shared/util"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQ struct {
	conn *amqp.Connection
}

func NewRabbitMQ(uri string) (*rabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	rmq := new(rabbitMQ)
	rmq.conn = conn

	return rmq, nil
}

func (r *rabbitMQ) Close() {
	if r.conn == nil {
		return
	}
	util.CloseAndLog(r.conn, "RabbitMQ")
}
