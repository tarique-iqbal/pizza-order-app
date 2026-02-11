package messaging

import (
	"encoding/json"
	"log"
	"restaurant-service/internal/shared/event"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewRabbitMQPublisher(amqpURL string) *RabbitMQPublisher {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		"email_exchange", // Exchange name
		"topic",          // Exchange type
		true,             // Durable
		false,            // Auto-delete
		false,            // Internal
		false,            // No-wait
		nil,              // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	return &RabbitMQPublisher{conn: conn, channel: ch}
}

func (p *RabbitMQPublisher) Publish(event event.Event) error {
	var body []byte
	var err error

	body, err = json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"email_exchange",     // Exchange name
		event.GetEventName(), // Routing key (e.g., user.registered, order.placed)
		false,                // Mandatory
		false,                // Immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v body: %s event: %s", err, body, event.GetEventName())
	}
	return err
}

func (p *RabbitMQPublisher) Close() {
	p.channel.Close()
	p.conn.Close()
}
