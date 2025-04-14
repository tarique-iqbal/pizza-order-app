package email_test

import (
	"errors"
	"testing"

	aEmail "email-service/internal/application/email"
	dEmail "email-service/internal/domain/email"

	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	Called bool
	Passed dEmail.EventPayload
	Err    error
}

func (m *mockHandler) Handle(event dEmail.EventPayload) error {
	m.Called = true
	m.Passed = event
	return m.Err
}

func TestDispatch_CallsRegisteredHandler(t *testing.T) {
	dispatcher := aEmail.NewEventDispatcher()
	mock := &mockHandler{}

	dispatcher.Register("user.registered", mock)

	event := dEmail.EventPayload{Name: "user.registered", Data: []byte(`{}`)}
	err := dispatcher.Dispatch(event)

	assert.NoError(t, err)
	assert.True(t, mock.Called)
	assert.Equal(t, event, mock.Passed)
}

func TestDispatch_NoHandler(t *testing.T) {
	dispatcher := aEmail.NewEventDispatcher()

	event := dEmail.EventPayload{Name: "user.unknown", Data: []byte(`{}`)}
	err := dispatcher.Dispatch(event)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no handler registered for event: user.unknown")
}

func TestDispatch_HandlerReturnsError(t *testing.T) {
	var dispatcher dEmail.EventDispatcher = aEmail.NewEventDispatcher()
	mock := &mockHandler{Err: errors.New("handler error")}

	dispatcher.Register("user.registered", mock)

	event := dEmail.EventPayload{Name: "user.registered", Data: []byte(`{}`)}
	err := dispatcher.Dispatch(event)

	assert.Error(t, err)
	assert.Equal(t, "handler error", err.Error())
}
