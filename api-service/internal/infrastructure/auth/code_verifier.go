package auth

import (
	"api-service/internal/domain/auth"
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

func (s *CodeVerificationService) Verify(email string, code string) error {
	emailVerification, err := s.emailVerificationRepo.FindByEmail(email)
	if err != nil {
		return auth.ErrCodeInvalid
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

	emailVerification.IsUsed = true
	if err := s.emailVerificationRepo.Updates(emailVerification); err != nil {
		return errors.New("failed to mark code as used")
	}

	return nil
}
