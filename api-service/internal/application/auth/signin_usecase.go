package auth

import (
	"api-service/internal/domain/auth"
	"api-service/internal/domain/user"
	"errors"
)

type SignInUseCase struct {
	repo   user.UserRepository
	hasher auth.PasswordHasher
	jwt    auth.JWTService
}

func NewSignInUseCase(
	repo user.UserRepository,
	hasher auth.PasswordHasher,
	jwt auth.JWTService,
) *SignInUseCase {
	return &SignInUseCase{repo: repo, hasher: hasher, jwt: jwt}
}

func (uc *SignInUseCase) Execute(email string, password string) (string, error) {
	user, err := uc.repo.FindByEmail(email)
	if user == nil {
		return "", errors.New("no record found")
	}

	if err != nil {
		return "", errors.New("internal server error")
	}

	status := uc.hasher.Compare(user.Password, password)
	if !status {
		return "", errors.New("invalid credentials")
	}

	token, err := uc.jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
