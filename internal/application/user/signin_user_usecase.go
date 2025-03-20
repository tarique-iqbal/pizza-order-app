package user

import (
	"errors"
	"pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/security"
)

type SignInUserUseCase interface {
	Execute(email string, password string) (string, error)
}

type signInUserUseCase struct {
	userRepo user.UserRepository
}

func NewSignInUserUseCase(repo user.UserRepository) SignInUserUseCase {
	return &signInUserUseCase{userRepo: repo}
}

func (uc *signInUserUseCase) Execute(email string, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if user == nil {
		return "", errors.New("no record found")
	}

	if err != nil {
		return "", errors.New("internal server error")
	}

	status := security.ComparePassword(user.Password, password)
	if !status {
		return "", errors.New("invalid credentials")
	}

	token, err := security.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
