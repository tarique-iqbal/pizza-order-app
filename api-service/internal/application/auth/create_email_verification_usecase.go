package auth

import (
	"api-service/internal/domain/auth"
	"api-service/internal/shared/event"
	"strings"
	"time"
)

type CreateEmailVerificationUseCase struct {
	repo      auth.EmailVerificationRepository
	otp       auth.OTPGenerator
	publisher event.EventPublisher
}

func NewCreateEmailVerificationUseCase(
	repo auth.EmailVerificationRepository,
	otp auth.OTPGenerator,
	publisher event.EventPublisher,
) *CreateEmailVerificationUseCase {
	return &CreateEmailVerificationUseCase{repo: repo, otp: otp, publisher: publisher}
}

func (uc *CreateEmailVerificationUseCase) Execute(input EmailVerificationRequestDTO) error {
	email := strings.ToLower(input.Email)

	code, err := uc.otp.Generate(true)
	if err != nil {
		return err
	}

	verification := &auth.EmailVerification{
		Email:     email,
		Code:      code,
		IsUsed:    false,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
	}

	existing, err := uc.repo.FindByEmail(email)
	if err != nil {
		return err
	}

	if existing == nil {
		if err := uc.repo.Create(verification); err != nil {
			return err
		}
	} else if existing.IsUsed {
		return nil
	} else {
		existing.Code = code
		existing.ExpiresAt = verification.ExpiresAt

		if err := uc.repo.Updates(existing); err != nil {
			return err
		}
	}

	event := EmailVerificationCreatedEvent{
		Email:     email,
		Code:      code,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	event.Name = event.GetEventName()

	return uc.publisher.Publish(event)
}
