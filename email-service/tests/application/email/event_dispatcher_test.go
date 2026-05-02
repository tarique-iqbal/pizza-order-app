package email_test

import (
	"errors"
	"testing"

	emailapp "email-service/internal/application/email"
	"email-service/internal/domain/email"

	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	Called bool
	Passed email.EventPayload
	Err    error
}

func (m *mockHandler) Handle(event email.EventPayload) error {
	m.Called = true
	m.Passed = event
	return m.Err
}

func TestDispatch_CallsRegisteredHandler(t *testing.T) {
	dispatcher := emailapp.NewEventDispatcher()
	mock := &mockHandler{}

	dispatcher.Register("user.registered", mock)

	event := email.EventPayload{Name: "user.registered", Data: []byte(`{}`)}
	err := dispatcher.Dispatch(event)

	assert.NoError(t, err)
	assert.True(t, mock.Called)
	assert.Equal(t, event, mock.Passed)
}

func TestDispatch_NoHandler(t *testing.T) {
	dispatcher := emailapp.NewEventDispatcher()

	event := email.EventPayload{Name: "user.unknown", Data: []byte(`{}`)}
	err := dispatcher.Dispatch(event)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no handler registered for event: user.unknown")
}

func TestDispatch_HandlerReturnsError(t *testing.T) {
	var dispatcher email.EventDispatcher = emailapp.NewEventDispatcher()
	mock := &mockHandler{Err: errors.New("handler error")}

	dispatcher.Register("user.registered", mock)

	event := email.EventPayload{Name: "user.registered", Data: []byte(`{}`)}
	err := dispatcher.Dispatch(event)

	assert.Error(t, err)
	assert.Equal(t, "handler error", err.Error())
}
