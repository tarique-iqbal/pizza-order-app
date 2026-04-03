package auth

import "context"

type RefreshTokenRepository interface {
	Save(ctx context.Context, hashedToken string, userID int, ttlSeconds int64) error
	Find(ctx context.Context, hashedToken string) (int, error)
	Delete(ctx context.Context, hashedToken string) error
}
