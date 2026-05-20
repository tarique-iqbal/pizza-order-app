package messaging_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"

	"restaurant-service/internal/domain/restaurant"
	"restaurant-service/internal/infrastructure/messaging"
)

type mockDispatcher struct {
	DispatchCalled bool
	EventReceived  restaurant.EventPayload
	Fail           bool
}

func (m *mockDispatcher) Register(eventName string, handler restaurant.EventHandler) {}

func (m *mockDispatcher) Dispatch(ctx context.Context, event restaurant.EventPayload) error {
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
	t.Skip("Skipping this test temporarily")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := &mockDispatcher{}
	msgs := make(chan amqp091.Delivery)

	consumer := &messaging.RabbitMQConsumer{}
	go consumer.Run(ctx, dispatcher)

	// Simulate a valid message
	msgs <- makeDelivery([]byte(`{"business_name": "Pizza Palace"}`), "restaurant.initiated", false)

	time.Sleep(100 * time.Millisecond) // Give goroutine time to process

	assert.True(t, dispatcher.DispatchCalled)
	assert.Equal(t, "restaurant.initiated", dispatcher.EventReceived.Name)
}

func TestRabbitMQConsumer_Run_DispatchFailsOnce(t *testing.T) {
	t.Skip("Skipping this test temporarily")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := &mockDispatcher{Fail: true}
	msgs := make(chan amqp091.Delivery)

	consumer := &messaging.RabbitMQConsumer{}
	go consumer.Run(ctx, dispatcher)

	msgs <- makeDelivery([]byte(`{}`), "restaurant.initiated", false)

	time.Sleep(100 * time.Millisecond)

	assert.True(t, dispatcher.DispatchCalled)
}
