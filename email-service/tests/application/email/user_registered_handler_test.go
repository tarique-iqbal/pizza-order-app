package email_test

import (
	"errors"
	"os"
	"testing"

	"email-service/internal/application/email"
	dEmail "email-service/internal/domain/email"

	"github.com/stretchr/testify/assert"
)

var (
	_ dEmail.Sender         = (*mockSender)(nil)
	_ dEmail.TemplateLoader = (*mockTemplateLoader)(nil)
)

type mockSender struct {
	To      string
	Subject string
	Body    string
	Err     error
}

func (m *mockSender) SendEmail(to, subject, body string) error {
	m.To = to
	m.Subject = subject
	m.Body = body
	return m.Err
}

type mockTemplateLoader struct {
	RenderCount int
	Fail        bool
}

func (m *mockTemplateLoader) Render(name string, data any) (string, error) {
	if m.Fail {
		return "", errors.New("template rendering failed")
	}
	m.RenderCount++
	switch name {
	case "user_welcome_email_subject.html":
		return "Welcome to MockApp!", nil
	case "user_welcome_email_body.html":
		return "Hello Alice, welcome to MockApp!", nil
	}
	return "", nil
}

func TestUserRegisteredHandler_Handle_Success(t *testing.T) {
	os.Setenv("APP_NAME", "MockApp")
	os.Setenv("SUPPORT_EMAIL", "support@mock.com")

	sender := &mockSender{}
	template := &mockTemplateLoader{}
	handler := email.NewUserRegisteredHandler(sender, template)

	event := dEmail.EventPayload{
		Name: "user.registered",
		Data: []byte(`{
			"email": "test@example.com",
			"first_name": "Alice",
			"role": "user",
			"timestamp": "2024-01-01T00:00:00Z"
		}`),
	}

	err := handler.Handle(event)

	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", sender.To)
	assert.Equal(t, "Welcome to MockApp!", sender.Subject)
	assert.Equal(t, "Hello Alice, welcome to MockApp!", sender.Body)
	assert.Equal(t, 2, template.RenderCount) // subject + body
}

func TestUserRegisteredHandler_Handle_InvalidJSON(t *testing.T) {
	sender := &mockSender{}
	template := &mockTemplateLoader{}
	handler := email.NewUserRegisteredHandler(sender, template)

	event := dEmail.EventPayload{
		Name: "user.registered",
		Data: []byte(`{invalid}`),
	}

	err := handler.Handle(event)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestUserRegisteredHandler_Handle_TemplateRenderFails(t *testing.T) {
	sender := &mockSender{}
	template := &mockTemplateLoader{Fail: true}
	handler := email.NewUserRegisteredHandler(sender, template)

	event := dEmail.EventPayload{
		Name: "user.registered",
		Data: []byte(`{
			"email": "test@example.com",
			"first_name": "Alice",
			"role": "owner",
			"timestamp": "2024-01-01T00:00:00Z"
		}`),
	}

	err := handler.Handle(event)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template rendering failed")
}

func TestUserRegisteredHandler_Handle_EmailSendFails(t *testing.T) {
	sender := &mockSender{Err: errors.New("smtp error")}
	template := &mockTemplateLoader{}
	handler := email.NewUserRegisteredHandler(sender, template)

	event := dEmail.EventPayload{
		Name: "user.registered",
		Data: []byte(`{
			"email": "test@example.com",
			"first_name": "Alice",
			"role": "user",
			"timestamp": "2024-01-01T00:00:00Z"
		}`),
	}

	err := handler.Handle(event)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "smtp error")
}
