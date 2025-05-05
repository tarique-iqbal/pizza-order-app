package auth

import "context"

type CodeVerifier interface {
	Verify(ctx context.Context, email string, code string) error
}
