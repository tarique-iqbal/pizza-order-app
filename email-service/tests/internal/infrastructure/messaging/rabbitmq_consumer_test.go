package messaging_test

import (
	"context"
	"email-service/internal/domain/email"
	"email-service/internal/infrastructure/messaging"
	"errors"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

type mockDispatcher struct {
	DispatchCalled bool
	EventReceived  email.EventPayload
	Fail           bool
}

func (m *mockDispatcher) Register(eventName string, handler email.EventHandler) {}

func (m *mockDispatcher) Dispatch(event email.EventPayload) error {
	m.DispatchCalled = true
	m.EventReceived = event
	if m.Fail {
		return errors.New("fail dispatch")
	}
	return nil
}

func makeDelivery(body []byte, routingKey string, redelivered bool) amqp091.Delivery {
	msg := amqp091.Delivery{
		Body:        body,
		RoutingKey:  routingKey,
		Redelivered: redelivered,
	}
	return msg
}

func TestRabbitMQConsumer_Run_DispatchSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := &mockDispatcher{}
	msgs := make(chan amqp091.Delivery)

	consumer := &messaging.RabbitMQConsumer{}
	go consumer.Run(ctx, msgs, dispatcher)

	// Simulate a valid message
	msgs <- makeDelivery([]byte(`{"email":"test@example.com"}`), "user.registered", false)

	time.Sleep(100 * time.Millisecond) // Give goroutine time to process

	assert.True(t, dispatcher.DispatchCalled)
	assert.Equal(t, "user.registered", dispatcher.EventReceived.Name)
}

func TestRabbitMQConsumer_Run_DispatchFailsOnce(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := &mockDispatcher{Fail: true}
	msgs := make(chan amqp091.Delivery)

	consumer := &messaging.RabbitMQConsumer{}
	go consumer.Run(ctx, msgs, dispatcher)

	msgs <- makeDelivery([]byte(`{}`), "user.registered", false)

	time.Sleep(100 * time.Millisecond)

	assert.True(t, dispatcher.DispatchCalled)
}
