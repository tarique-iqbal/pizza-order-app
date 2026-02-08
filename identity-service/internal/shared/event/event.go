package event

type Event interface {
	GetEventName() string
}

type EventPublisher interface {
	Publish(event Event) error
}
