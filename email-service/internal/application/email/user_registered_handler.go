package email

import (
	"email-service/internal/domain/email"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
		Role      string `json:"role"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.Unmarshal(event.Data, &payload); err != nil {
		return err
	}

	validRoles := map[string]bool{"User": true, "Owner": true}
	if !validRoles[payload.Role] {
		return fmt.Errorf("invalid role in event payload: %s", payload.Role)
	}

	subjectTemplate := fmt.Sprintf("%s_welcome_email_subject.html", strings.ToLower(payload.Role))
	bodyTemplate := fmt.Sprintf("%s_welcome_email_body.html", strings.ToLower(payload.Role))

	appName := os.Getenv("APP_NAME")
	supportEmail := os.Getenv("SUPPORT_EMAIL")

	subject, err := h.template.Render(subjectTemplate, map[string]string{
		"app_name": appName,
	})
	if err != nil {
		return err
	}

	body, err := h.template.Render(bodyTemplate, map[string]string{
		"first_name":    payload.FirstName,
		"app_name":      appName,
		"support_email": supportEmail,
	})
	if err != nil {
		return err
	}

	return h.sender.SendEmail(payload.Email, subject, body)
}
