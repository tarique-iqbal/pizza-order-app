package auth

type EmailVerificationCreatedEvent struct {
	Email     string `json:"email"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
}

func (e EmailVerificationCreatedEvent) GetEventName() string {
	return "email.verification_created"
}
