package auth_test

import (
	"context"
	"identity-service/internal/application/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	"identity-service/tests/infrastructure/db/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

var login *auth.Login

func setupLogin(t *testing.T) *auth.Login {
	ts := testStorage()
	flushRedis(t, ts.Redis)
	truncateTables(ts.DB)

	if err := fixtures.LoadUserFixtures(ts.DB); err != nil {
		panic(err)
	}

	repo := persistence.NewUserRepository(ts.DB)
	hasher := security.NewPasswordHasher()
	jwt := security.NewJWTManager("TestSecretKey")
	refreshTokenRepo := persistence.NewRefreshTokenRepository(ts.Redis)
	refreshTokenManager := security.NewRefreshTokenManager()

	return auth.NewLogin(repo, hasher, jwt, refreshTokenRepo, refreshTokenManager)
}

func TestLogin_Success(t *testing.T) {
	login := setupLogin(t)

	response, err := login.Execute(context.Background(), "john.doe@example.com", "plainPassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
}

func TestLogin_InvalidPassword(t *testing.T) {
	login := setupLogin(t)

	response, err := login.Execute(context.Background(), "john.doe@example.com", "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, response.AccessToken)
	assert.Empty(t, response.RefreshToken)
}

func TestLogin_UserNotFound(t *testing.T) {
	login := setupLogin(t)

	response, err := login.Execute(context.Background(), "notfound@example.com", "password")
	assert.Error(t, err)
	assert.Empty(t, response.AccessToken)
	assert.Empty(t, response.RefreshToken)
}
