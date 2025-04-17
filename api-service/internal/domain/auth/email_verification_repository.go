package auth

type EmailVerificationRepository interface {
	Create(emailVerification *EmailVerification) error
	Updates(emailVerification *EmailVerification) error
	FindByEmail(email string) (*EmailVerification, error)
}
