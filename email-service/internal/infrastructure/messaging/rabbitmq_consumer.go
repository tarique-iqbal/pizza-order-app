package messaging

import (
	"context"
	"log"
	"time"

	"email-service/internal/domain/email"

	"github.com/rabbitmq/amqp091-go"
)

type deliveryChan <-chan amqp091.Delivery
type dispatcher email.EventDispatcher

type RabbitMQConsumer struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewRabbitMQConsumer(amqpURL string) *RabbitMQConsumer {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		"email_queue", // Queue name
		true,          // Durable
		false,         // Auto-delete
		false,         // Exclusive
		false,         // No-wait
		nil,           // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	routingKeys := []string{
		"email.verification_created",
		"user.registered",
	}
	for _, key := range routingKeys {
		err = ch.QueueBind(
			q.Name,           // Queue name
			key,              // Routing key
			"email_exchange", // Exchange name
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to bind queue: %v", err)
		}
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	return &RabbitMQConsumer{conn: conn, channel: ch}
}

func (c *RabbitMQConsumer) GetMessages() deliveryChan {
	msgs, err := c.channel.Consume(
		"email_queue",          // Queue name
		"email_queue_consumer", // Consumer name
		false,                  // Auto-Ack Disabled
		false,                  // Exclusive
		false,                  // No-local
		false,                  // No-wait
		nil,                    // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	return msgs
}

func (c *RabbitMQConsumer) Run(ctx context.Context, msgs deliveryChan, dispatcher dispatcher) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer stopped gracefully")
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Println("Message channel closed")
				return
			}

			go func(msg amqp091.Delivery) {
				event := email.EventPayload{
					Name: msg.RoutingKey,
					Data: msg.Body,
				}
				if err := dispatcher.Dispatch(event); err != nil {
					log.Println("Error dispatching:", err)

					if msg.Redelivered {
						log.Println("Message already redelivered, rejecting permanently.")
						_ = msg.Nack(false, false) // Do not requeue
					} else {
						log.Println("Requeueing message for retry...")
						time.Sleep(2 * time.Second)
						_ = msg.Nack(false, true) // Requeue the message
					}
					return
				}

				msg.Ack(false)
			}(msg)
		}
	}
}

func (c *RabbitMQConsumer) Close() {
	c.channel.Close()
	c.conn.Close()
}
