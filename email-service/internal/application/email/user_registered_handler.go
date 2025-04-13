package email

import (
	"email-service/internal/domain/email"
	"encoding/json"
	"os"
)

type UserRegisteredHandler struct {
	sender   email.Sender
	template email.TemplateLoader
}

func NewUserRegisteredHandler(sender email.Sender, template email.TemplateLoader) *UserRegisteredHandler {
	return &UserRegisteredHandler{sender: sender, template: template}
}

func (h *UserRegisteredHandler) Handle(event email.EventPayload) error {
	var payload struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.Unmarshal(event.Data, &payload); err != nil {
		return err
	}

	appName := os.Getenv("APP_NAME")
	supportEmail := os.Getenv("SUPPORT_EMAIL")
	subject, err := h.template.Render("welcome_email_subject.html", map[string]string{
		"app_name": appName,
	})
	if err != nil {
		return err
	}
	body, err := h.template.Render("welcome_email_body.html", map[string]string{
		"first_name":    payload.FirstName,
		"app_name":      appName,
		"support_email": supportEmail,
	})
	if err != nil {
		return err
	}

	return h.sender.SendEmail(payload.Email, subject, body)
}
