package user

import (
	"context"
	"identity-service/internal/domain/auth"
	"identity-service/internal/domain/user"
	"identity-service/internal/shared/event"
	"log"
	"time"
)

const defaultStatus = "active"

type Register struct {
	emailVerifier auth.EmailVerifier
	repo          user.UserRepository
	hasher        auth.PasswordHasher
	publisher     event.EventPublisher
}

func NewRegister(
	emailVerifier auth.EmailVerifier,
	repo user.UserRepository,
	hasher auth.PasswordHasher,
	publisher event.EventPublisher,
) *Register {
	return &Register{
		emailVerifier: emailVerifier,
		repo:          repo,
		hasher:        hasher,
		publisher:     publisher,
	}
}

func (uc *Register) Execute(ctx context.Context, input RegisterRequest) (Response, error) {
	if err := uc.emailVerifier.Verify(ctx, input.Email, input.Code); err != nil {
		return Response{}, err
	}

	hashedPassword, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return Response{}, err
	}

	newUser := user.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
		Role:      input.Role,
		Status:    defaultStatus,
	}

	if err := uc.repo.Create(ctx, &newUser); err != nil {
		return Response{}, err
	}

	event := UserRegistered{
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		Role:      newUser.Role,
		Timestamp: newUser.CreatedAt.Format(time.RFC3339),
	}
	event.EventName = event.GetEventName()

	if err := uc.publisher.Publish(event); err != nil {
		log.Println("Failed to publish user.registered event:", err)
	}

	response := Response{
		ID:        newUser.ID,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		Role:      newUser.Role,
	}

	return response, nil
}
