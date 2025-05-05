package auth

import (
	"api-service/internal/domain/auth"
	"context"
	"errors"
	"time"
)

type CodeVerificationService struct {
	emailVerificationRepo auth.EmailVerificationRepository
}

func NewCodeVerificationService(
	emailVerificationRepo auth.EmailVerificationRepository,
) auth.CodeVerifier {
	return &CodeVerificationService{emailVerificationRepo: emailVerificationRepo}
}

func (s *CodeVerificationService) Verify(ctx context.Context, email string, code string) error {
	emailVerification, err := s.emailVerificationRepo.FindByEmail(ctx, email)
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
	if err := s.emailVerificationRepo.Updates(ctx, emailVerification); err != nil {
		return errors.New("failed to mark code as used")
	}

	return nil
}
