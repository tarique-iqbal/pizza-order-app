package user

import (
	"log"
	"pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/security"
	"pizza-order-api/internal/shared/event"
	"time"
)

const defaultRole = "user"
const defaultStatus = "Active"
const defaultVerified = "No"

type CreateUserUseCase struct {
	repo      user.UserRepository
	publisher event.EventPublisher
}

func NewCreateUserUseCase(repo user.UserRepository, publisher event.EventPublisher) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo, publisher: publisher}
}

func (uc *CreateUserUseCase) Execute(input UserCreateDTO) (UserResponseDTO, error) {
	hashedPassword, err := security.HashPassword(input.Password)
	if err != nil {
		return UserResponseDTO{}, err
	}

	newUser := user.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
		Role:      defaultRole,
		Status:    defaultStatus,
		Verified:  defaultVerified,
		LoggedAt:  nil,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	if err := uc.repo.Create(&newUser); err != nil {
		return UserResponseDTO{}, err
	}

	event := UserCreatedEvent{
		UserID:    newUser.ID,
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Timestamp: newUser.CreatedAt.Format(time.RFC3339),
	}

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
