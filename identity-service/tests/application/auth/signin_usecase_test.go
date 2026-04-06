package auth_test

import (
	"context"
	"identity-service/internal/application/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	"identity-service/tests/infrastructure/db/fixtures"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var signInUC *auth.SignInUseCase

func setupSignInUseCase(t *testing.T) *auth.SignInUseCase {
	ts := testStorage()
	flushRedis(t, ts.Redis)
	truncateTables(ts.DB)

	if err := fixtures.LoadUserFixtures(ts.DB); err != nil {
		panic(err)
	}

	repo := persistence.NewUserRepository(ts.DB)
	hasher := security.NewPasswordHasher()
	jwt := security.NewJWTService("TestSecretKey")
	refreshTokenRepo := persistence.NewRefreshTokenRepository(ts.Redis)
	refreshTokenService := security.NewRefreshTokenService()

	return auth.NewSignInUseCase(repo, hasher, jwt, refreshTokenRepo, refreshTokenService)
}

func flushRedis(t *testing.T, client *redis.Client) {
	err := client.FlushDB(context.Background()).Err()
	require.NoError(t, err)
}

func TestSignInUseCase_Success(t *testing.T) {
	signInUC := setupSignInUseCase(t)

	response, err := signInUC.Execute(context.Background(), "john.doe@example.com", "plainPassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
}

func TestSignInUseCase_InvalidPassword(t *testing.T) {
	signInUC := setupSignInUseCase(t)

	response, err := signInUC.Execute(context.Background(), "john.doe@example.com", "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, response.AccessToken)
	assert.Empty(t, response.RefreshToken)
}

func TestSignInUseCase_UserNotFound(t *testing.T) {
	signInUC := setupSignInUseCase(t)

	response, err := signInUC.Execute(context.Background(), "notfound@example.com", "password")
	assert.Error(t, err)
	assert.Empty(t, response.AccessToken)
	assert.Empty(t, response.RefreshToken)
}
