package auth

import (
	"context"
	"errors"
	"identity-service/internal/domain/auth"
	"identity-service/internal/domain/user"
)

const refreshTokenExpiry = 7

type Login struct {
	userRepo            user.UserRepository
	passwordHasher      auth.PasswordHasher
	jwtManager          auth.JWTManager
	refreshTokenRepo    auth.RefreshTokenRepository
	refreshTokenManager auth.RefreshTokenManager
}

func NewLogin(
	userRepo user.UserRepository,
	passwordHasher auth.PasswordHasher,
	jwtManager auth.JWTManager,
	refreshTokenRepo auth.RefreshTokenRepository,
	refreshTokenManager auth.RefreshTokenManager,
) *Login {
	return &Login{
		userRepo:            userRepo,
		passwordHasher:      passwordHasher,
		jwtManager:          jwtManager,
		refreshTokenRepo:    refreshTokenRepo,
		refreshTokenManager: refreshTokenManager,
	}
}

func (uc *Login) Execute(
	ctx context.Context,
	email string,
	password string,
) (TokenResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if user == nil {
		return TokenResponse{}, errors.New("no record found")
	}

	if err != nil {
		return TokenResponse{}, errors.New("internal server error")
	}

	status := uc.passwordHasher.Compare(user.Password, password)
	if !status {
		return TokenResponse{}, errors.New("invalid credentials")
	}

	accessToken, err := uc.jwtManager.Generate(user.ID, user.Role)
	if err != nil {
		return TokenResponse{}, err
	}

	refreshToken, err := uc.refreshTokenManager.Generate()
	if err != nil {
		return TokenResponse{}, err
	}

	hashedToken, _ := uc.refreshTokenManager.Hash(refreshToken)

	const ttlSeconds = int64(refreshTokenExpiry * 24 * 3600)
	err = uc.refreshTokenRepo.Save(ctx, hashedToken, user.ID, ttlSeconds)
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
