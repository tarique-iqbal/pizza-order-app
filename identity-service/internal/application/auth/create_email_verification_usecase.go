package auth

import (
	"context"
	"identity-service/internal/domain/auth"
	"identity-service/internal/shared/event"
	"strings"
	"time"
)

const tokenExpiryMinutes = 15

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

func (uc *CreateEmailVerificationUseCase) Execute(
	ctx context.Context,
	input EmailVerificationRequestDTO,
) error {
	email := strings.ToLower(input.Email)

	code, err := uc.otp.Generate(true)
	if err != nil {
		return err
	}

	verification := &auth.EmailVerification{
		Email:     email,
		Code:      code,
		IsUsed:    false,
		ExpiresAt: time.Now().Add(time.Duration(tokenExpiryMinutes) * time.Minute),
	}

	existing, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if existing == nil {
		if err := uc.repo.Create(ctx, verification); err != nil {
			return err
		}
	} else if existing.IsUsed {
		return nil
	} else {
		existing.Code = code
		existing.ExpiresAt = verification.ExpiresAt

		if err := uc.repo.Updates(ctx, existing); err != nil {
			return err
		}
	}

	event := EmailVerificationCreatedEvent{
		Email:     email,
		Code:      code,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	event.EventName = event.GetEventName()

	return uc.publisher.Publish(event)
}
