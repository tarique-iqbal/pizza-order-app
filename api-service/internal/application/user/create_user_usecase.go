package user

import (
	"api-service/internal/domain/auth"
	"api-service/internal/domain/user"
	"api-service/internal/shared/event"
	"context"
	"log"
	"time"
)

const defaultStatus = "active"

type CreateUserUseCase struct {
	codeVerifier auth.CodeVerifier
	repo         user.UserRepository
	hasher       auth.PasswordHasher
	publisher    event.EventPublisher
}

func NewCreateUserUseCase(
	codeVerifier auth.CodeVerifier,
	repo user.UserRepository,
	hasher auth.PasswordHasher,
	publisher event.EventPublisher,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		codeVerifier: codeVerifier,
		repo:         repo,
		hasher:       hasher,
		publisher:    publisher,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input UserCreateDTO) (UserResponseDTO, error) {
	if err := uc.codeVerifier.Verify(ctx, input.Email, input.Code); err != nil {
		return UserResponseDTO{}, err
	}

	hashedPassword, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return UserResponseDTO{}, err
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
		return UserResponseDTO{}, err
	}

	event := UserCreatedEvent{
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		Role:      newUser.Role,
		Timestamp: newUser.CreatedAt.Format(time.RFC3339),
	}
	event.EventName = event.GetEventName()

	if err := uc.publisher.Publish(event); err != nil {
		log.Println("Failed to publish user.registered event:", err)
	}

	response := UserResponseDTO{
		ID:        newUser.ID,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		Role:      newUser.Role,
	}

	return response, nil
}
