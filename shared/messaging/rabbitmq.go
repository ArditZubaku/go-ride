// Package messaging provides messaging system capabilities
// It currently uses RabbitMQ
package messaging

import (
	"context"
	"fmt"
	"log"

	"ride-sharing/shared/util"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		util.CloseOrLog(conn, "RabbitMQ connection")
		return nil, fmt.Errorf("failed to create channel on RabbitMQ: %v", err)
	}

	rmq := new(RabbitMQ)
	rmq.conn = conn
	rmq.ch = ch

	// Monitor connection closures
	go rmq.monitorConnection()

	// Monitor channel closures
	go rmq.monitorChannel()

	if err := rmq.setupExchanges(); err != nil {
		rmq.Close()
		return nil, fmt.Errorf("failed to setup exchanges on RabbitMQ: %v", err)
	}

	if err := rmq.setupQueues(); err != nil {
		rmq.Close()
		return nil, fmt.Errorf("failed to setup queues on RabbitMQ: %v", err)
	}

	return rmq, nil
}

// monitorConnection listens for connection close events and logs them
func (r *RabbitMQ) monitorConnection() {
	closeChan := r.conn.NotifyClose(make(chan *amqp.Error))
	err := <-closeChan
	if err != nil {
		log.Printf("[RabbitMQ] Connection lost: %v (code: %d, reason: %s)", err, err.Code, err.Reason)
	}
	// If err is nil, connection was closed cleanly (e.g., via Close()), no need to log
}

// monitorChannel listens for channel close events and logs them
func (r *RabbitMQ) monitorChannel() {
	closeChan := r.ch.NotifyClose(make(chan *amqp.Error))
	err := <-closeChan
	if err != nil {
		log.Printf("[RabbitMQ] Channel lost: %v (code: %d, reason: %s)", err, err.Code, err.Reason)
	}
	// If err is nil, channel was closed cleanly (e.g., via Close()), no need to log
}

func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.ch
}

func (r *RabbitMQ) Publish(
	ctx context.Context,
	routingKey string,
	message string,
) error {
	// I am passing the ctx just because - but this API doesn't honour it...
	return r.ch.PublishWithContext(
		ctx,
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (r *RabbitMQ) Close() {
	if r.conn == nil {
		return
	}
	util.CloseOrLog(r.conn, "RabbitMQ connection")

	if r.ch == nil {
		return
	}
	util.CloseOrLog(r.ch, "RabbitMQ channel")
}

func (r *RabbitMQ) setupExchanges() error {
	return nil
}

func (r *RabbitMQ) setupQueues() error {
	_, err := r.ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	return err
}
