package auth

type EmailVerificationCreatedEvent struct {
	Email     string `json:"email"`
	Code      string `json:"code"`
	EventName string `json:"event_name"`
	Timestamp string `json:"timestamp"`
}

func (e EmailVerificationCreatedEvent) GetEventName() string {
	return "email.verification_created"
}
