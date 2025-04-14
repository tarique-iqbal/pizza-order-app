package email

import (
	"email-service/internal/domain/email"
	"fmt"
)

type EventDispatcher struct {
	handlers map[string]email.EventHandler
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string]email.EventHandler),
	}
}

func (d *EventDispatcher) Register(eventName string, handler email.EventHandler) {
	d.handlers[eventName] = handler
}

func (d *EventDispatcher) Dispatch(event email.EventPayload) error {
	handler, ok := d.handlers[event.Name]
	if !ok {
		return fmt.Errorf("no handler registered for event: %s", event.Name)
	}
	return handler.Handle(event)
}
