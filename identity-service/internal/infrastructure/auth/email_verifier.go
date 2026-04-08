package auth

import (
	"context"
	"errors"
	"identity-service/internal/domain/auth"
	"time"
)

type emailVerifier struct {
	repo auth.EmailVerificationRepository
}

func NewEmailVerifier(
	repo auth.EmailVerificationRepository,
) auth.EmailVerifier {
	return &emailVerifier{repo: repo}
}

func (s *emailVerifier) Verify(ctx context.Context, email string, code string) error {
	emailVerification, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return auth.ErrCodeInvalid
	}

	if emailVerification == nil {
		return auth.ErrCodeNotIssued
	}

	if emailVerification.IsUsed {
		return auth.ErrCodeUsed
	}

	if emailVerification.Code != code {
		return auth.ErrCodeInvalid
	}

	if time.Now().After(emailVerification.ExpiresAt) {
		return auth.ErrCodeExpired
	}

	emailVerification.Code = "..."
	emailVerification.IsUsed = true
	if err := s.repo.Updates(ctx, emailVerification); err != nil {
		return errors.New("failed to mark code as used")
	}

	return nil
}
