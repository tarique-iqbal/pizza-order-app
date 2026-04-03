package auth

import (
	"context"
	"errors"
	"identity-service/internal/domain/auth"
	"identity-service/internal/domain/user"
)

type SignInUseCase struct {
	repo         user.UserRepository
	hasher       auth.PasswordHasher
	jwt          auth.JWTService
	refreshRepo  auth.RefreshTokenRepository
	refreshToken auth.RefreshTokenService
}

func NewSignInUseCase(
	repo user.UserRepository,
	hasher auth.PasswordHasher,
	jwt auth.JWTService,
	refreshRepo auth.RefreshTokenRepository,
	refreshToken auth.RefreshTokenService,
) *SignInUseCase {
	return &SignInUseCase{
		repo:         repo,
		hasher:       hasher,
		jwt:          jwt,
		refreshRepo:  refreshRepo,
		refreshToken: refreshToken,
	}
}

func (uc *SignInUseCase) Execute(
	ctx context.Context,
	email string,
	password string,
) (string, string, error) {
	user, err := uc.repo.FindByEmail(ctx, email)
	if user == nil {
		return "", "", errors.New("no record found")
	}

	if err != nil {
		return "", "", errors.New("internal server error")
	}

	status := uc.hasher.Compare(user.Password, password)
	if !status {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := uc.jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uc.refreshToken.Generate()
	if err != nil {
		return "", "", err
	}

	hashedToken, _ := uc.refreshToken.Hash(refreshToken)

	const ttlSeconds = int64(7 * 24 * 3600)
	err = uc.refreshRepo.Save(ctx, hashedToken, int(user.ID), ttlSeconds)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
