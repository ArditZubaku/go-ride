// Package messaging provides messaging system capabilities
// It currently uses RabbitMQ
package messaging

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"ride-sharing/shared/retry"
	"ride-sharing/shared/util"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	uri      string
	conn     *amqp.Connection
	ch       *amqp.Channel
	mu       sync.RWMutex
	shutdown bool
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	rmq := &RabbitMQ{
		uri:      uri,
		shutdown: false,
	}

	if err := rmq.connect(); err != nil {
		return nil, err
	}

	// Monitor connection closures
	go rmq.monitorConnection()

	// Monitor channel closures
	go rmq.monitorChannel()

	return rmq, nil
}

// connect establishes a new connection and channel, and sets up exchanges/queues
func (r *RabbitMQ) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	conn, err := amqp.Dial(r.uri)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		util.CloseOrLog(conn, "RabbitMQ connection")
		return fmt.Errorf("failed to create channel on RabbitMQ: %v", err)
	}

	r.conn = conn
	r.ch = ch

	if err := r.setupExchanges(); err != nil {
		util.CloseOrLog(ch, "RabbitMQ channel")
		util.CloseOrLog(conn, "RabbitMQ connection")
		return fmt.Errorf("failed to setup exchanges on RabbitMQ: %v", err)
	}

	if err := r.setupQueues(); err != nil {
		util.CloseOrLog(ch, "RabbitMQ channel")
		util.CloseOrLog(conn, "RabbitMQ connection")
		return fmt.Errorf("failed to setup queues on RabbitMQ: %v", err)
	}

	return nil
}

// monitorConnection listens for connection close events and attempts to reconnect
func (r *RabbitMQ) monitorConnection() {
	for {
		r.mu.RLock()
		if r.shutdown || r.conn == nil {
			r.mu.RUnlock()
			return
		}
		closeChan := r.conn.NotifyClose(make(chan *amqp.Error))
		r.mu.RUnlock()

		err := <-closeChan
		if err == nil {
			// Connection closed cleanly (e.g., via Close())
			return
		}

		log.Printf("[RabbitMQ] Connection lost: %v (code: %d, reason: %s). Attempting to reconnect...", err, err.Code, err.Reason)

		// Attempt to reconnect with retry
		ctx := context.Background()
		retryCfg := retry.Config{
			MaxRetries:  10, // Try up to 10 times
			InitialWait: 1 * time.Second,
			MaxWait:     30 * time.Second,
		}

		reconnectErr := retry.WithBackoff(ctx, retryCfg, func() error {
			if r.shutdown {
				return fmt.Errorf("shutdown in progress")
			}
			return r.connect()
		})

		if reconnectErr != nil {
			if r.shutdown {
				log.Printf("[RabbitMQ] Shutdown in progress, stopping reconnection attempts")
				return
			}
			log.Printf("[RabbitMQ] Failed to reconnect after retries: %v", reconnectErr)
			// Continue monitoring in case we want to retry again later
			time.Sleep(5 * time.Second)
		} else {
			log.Printf("[RabbitMQ] Successfully reconnected")
			// Restart channel monitoring for the new connection
			go r.monitorChannel()
		}
	}
}

// monitorChannel listens for channel close events and attempts to recreate the channel
// It loops to monitor each new channel after recreation
func (r *RabbitMQ) monitorChannel() {
	for {
		// Get the current channel's close notification channel
		r.mu.RLock()
		if r.shutdown {
			r.mu.RUnlock()
			return
		}
		if r.conn == nil || r.ch == nil {
			r.mu.RUnlock()
			// Connection or channel is gone, exit (connection monitor will handle reconnection)
			return
		}
		closeChan := r.ch.NotifyClose(make(chan *amqp.Error))
		r.mu.RUnlock()

		// Wait for channel to close (this channel fires once per channel)
		err := <-closeChan
		if err == nil {
			// Channel closed cleanly (e.g., via Close())
			return
		}

		log.Printf("[RabbitMQ] Channel lost: %v (code: %d, reason: %s). Attempting to recreate...", err, err.Code, err.Reason)

		// Attempt to recreate channel with retry
		ctx := context.Background()
		retryCfg := retry.Config{
			MaxRetries:  5,
			InitialWait: 1 * time.Second,
			MaxWait:     10 * time.Second,
		}

		recreateErr := retry.WithBackoff(ctx, retryCfg, func() error {
			if r.shutdown {
				return fmt.Errorf("shutdown in progress")
			}

			r.mu.Lock()
			defer r.mu.Unlock()

			if r.conn == nil {
				return fmt.Errorf("connection is nil, cannot recreate channel")
			}

			ch, err := r.conn.Channel()
			if err != nil {
				return fmt.Errorf("failed to create channel: %v", err)
			}

			r.ch = ch

			// Re-setup queues after channel recreation
			if err := r.setupQueues(); err != nil {
				util.CloseOrLog(ch, "RabbitMQ channel")
				return fmt.Errorf("failed to setup queues: %v", err)
			}

			return nil
		})

		if recreateErr != nil {
			if r.shutdown {
				log.Printf("[RabbitMQ] Shutdown in progress, stopping channel recreation attempts")
				return
			}
			log.Printf("[RabbitMQ] Failed to recreate channel after retries: %v", recreateErr)
			// If channel recreation fails, connection might be lost too, so exit
			// The connection monitor will handle reconnection
			return
		}

		log.Printf("[RabbitMQ] Successfully recreated channel")
		// Loop back to monitor the NEW channel (NotifyClose only works once per channel)
	}
}

func (r *RabbitMQ) GetChannel() *amqp.Channel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ch
}

func (r *RabbitMQ) Publish(
	ctx context.Context,
	routingKey string,
	message string,
) error {
	retryCfg := retry.Config{
		MaxRetries:  3,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     2 * time.Second,
	}

	return retry.WithBackoff(ctx, retryCfg, func() error {
		r.mu.RLock()
		if r.shutdown {
			r.mu.RUnlock()
			return fmt.Errorf("RabbitMQ client is shutdown")
		}
		if r.ch == nil {
			r.mu.RUnlock()
			return fmt.Errorf("RabbitMQ channel is nil")
		}
		ch := r.ch
		r.mu.RUnlock()

		err := ch.PublishWithContext(
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

		// If publish fails due to connection/channel issues, the monitors will handle reconnection
		// and the next retry should succeed
		return err
	})
}

func (r *RabbitMQ) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.shutdown = true

	if r.ch != nil {
		util.CloseOrLog(r.ch, "RabbitMQ channel")
		r.ch = nil
	}

	if r.conn != nil {
		util.CloseOrLog(r.conn, "RabbitMQ connection")
		r.conn = nil
	}
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
