package email

import (
	"email-service/internal/domain/email"
	"encoding/json"
	"os"
)

type EmailVerificationCreatedHandler struct {
	sender   email.Sender
	template email.TemplateLoader
}

func NewEmailVerificationCreatedHandler(sender email.Sender, template email.TemplateLoader) *EmailVerificationCreatedHandler {
	return &EmailVerificationCreatedHandler{sender: sender, template: template}
}

func (h *EmailVerificationCreatedHandler) Handle(event email.EventPayload) error {
	var payload struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.Unmarshal(event.Data, &payload); err != nil {
		return err
	}

	appName := os.Getenv("APP_NAME")
	tokenExpiryMinutes := os.Getenv("TOKEN_EXPIRY_MINUTES")

	subject, err := h.template.Render("email_verification_subject.html", map[string]string{
		"app_name": appName,
	})
	if err != nil {
		return err
	}
	body, err := h.template.Render("email_verification_body.html", map[string]string{
		"email":                payload.Email,
		"code":                 payload.Code,
		"app_name":             appName,
		"token_expiry_minutes": tokenExpiryMinutes,
	})
	if err != nil {
		return err
	}

	return h.sender.SendEmail(payload.Email, subject, body)
}
