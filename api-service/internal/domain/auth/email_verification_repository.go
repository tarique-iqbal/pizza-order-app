package auth

import "context"

type EmailVerificationRepository interface {
	Create(ctx context.Context, emailVerification *EmailVerification) error
	Updates(ctx context.Context, emailVerification *EmailVerification) error
	FindByEmail(ctx context.Context, email string) (*EmailVerification, error)
}
