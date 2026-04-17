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
	input LoginRequest,
) (TokenResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if user == nil {
		return TokenResponse{}, errors.New("no record found")
	}

	if err != nil {
		return TokenResponse{}, errors.New("internal server error")
	}

	status := uc.passwordHasher.Compare(user.Password, input.Password)
	if !status {
		return TokenResponse{}, errors.New("invalid credentials")
	}

	accessToken, err := uc.jwtManager.Generate(user.ID.String(), user.Role)
	if err != nil {
		return TokenResponse{}, err
	}

	refreshToken, err := uc.refreshTokenManager.Generate()
	if err != nil {
		return TokenResponse{}, err
	}

	hashedToken := uc.refreshTokenManager.Hash(refreshToken)

	claims := auth.UserClaims{
		UserID: user.ID.String(),
		Role:   user.Role,
	}

	ttlSeconds := int64(refreshTokenExpiry) * 24 * 3600

	err = uc.refreshTokenRepo.Save(ctx, hashedToken, claims, ttlSeconds)
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
