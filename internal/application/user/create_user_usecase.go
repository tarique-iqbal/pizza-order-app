package user

import (
	"pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/security"
	"time"
)

const defaultRole = "user"

type CreateUserUseCase interface {
	Execute(UserCreateDTO) (UserResponseDTO, error)
}

type createUserUseCase struct {
	repo user.UserRepository
}

func NewCreateUserUseCase(repo user.UserRepository) CreateUserUseCase {
	return &createUserUseCase{repo}
}

func (uc *createUserUseCase) Execute(input UserCreateDTO) (UserResponseDTO, error) {
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
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	if err := uc.repo.Create(&newUser); err != nil {
		return UserResponseDTO{}, err
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
