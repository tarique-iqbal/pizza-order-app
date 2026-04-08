package auth

import "context"

type EmailVerificationRepository interface {
	Create(ctx context.Context, emailVerification *EmailVerification) error
	Updates(ctx context.Context, emailVerification *EmailVerification) error
	FindByEmail(ctx context.Context, email string) (*EmailVerification, error)
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, hashedToken string, userID int, ttlSeconds int64) error
	Find(ctx context.Context, hashedToken string) (int, error)
	Delete(ctx context.Context, hashedToken string) error
}
