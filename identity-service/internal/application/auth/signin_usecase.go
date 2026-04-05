package auth

import (
	"context"
	"errors"
	"identity-service/internal/domain/auth"
	"identity-service/internal/domain/user"
)

type SignInUseCase struct {
	userRepo            user.UserRepository
	passwordHasher      auth.PasswordHasher
	jwtService          auth.JWTService
	refreshTokenRepo    auth.RefreshTokenRepository
	refreshTokenService auth.RefreshTokenService
}

func NewSignInUseCase(
	userRepo user.UserRepository,
	passwordHasher auth.PasswordHasher,
	jwtService auth.JWTService,
	refreshTokenRepo auth.RefreshTokenRepository,
	refreshTokenService auth.RefreshTokenService,
) *SignInUseCase {
	return &SignInUseCase{
		userRepo:            userRepo,
		passwordHasher:      passwordHasher,
		jwtService:          jwtService,
		refreshTokenRepo:    refreshTokenRepo,
		refreshTokenService: refreshTokenService,
	}
}

func (uc *SignInUseCase) Execute(
	ctx context.Context,
	email string,
	password string,
) (string, string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if user == nil {
		return "", "", errors.New("no record found")
	}

	if err != nil {
		return "", "", errors.New("internal server error")
	}

	status := uc.passwordHasher.Compare(user.Password, password)
	if !status {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := uc.jwtService.Generate(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uc.refreshTokenService.Generate()
	if err != nil {
		return "", "", err
	}

	hashedToken, _ := uc.refreshTokenService.Hash(refreshToken)

	const ttlSeconds = int64(7 * 24 * 3600)
	err = uc.refreshTokenRepo.Save(ctx, hashedToken, int(user.ID), ttlSeconds)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
