package user

import (
	"api-service/internal/domain/auth"
	"api-service/internal/domain/user"
	"api-service/internal/shared/event"
	"log"
	"time"
)

const defaultStatus = "Active"

type CreateUserUseCase struct {
	repo      user.UserRepository
	hasher    auth.PasswordHasher
	publisher event.EventPublisher
}

func NewCreateUserUseCase(
	repo user.UserRepository,
	hasher auth.PasswordHasher,
	publisher event.EventPublisher,
) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo, hasher: hasher, publisher: publisher}
}

func (uc *CreateUserUseCase) Execute(input UserCreateDTO) (UserResponseDTO, error) {
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
		LoggedAt:  nil,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	if err := uc.repo.Create(&newUser); err != nil {
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
