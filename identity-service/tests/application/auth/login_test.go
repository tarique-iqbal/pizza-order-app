package auth_test

import (
	"context"
	"identity-service/internal/application/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	"identity-service/tests/infrastructure/db/fixtures"
	"identity-service/tests/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var login *auth.Login

func setupLogin(t *testing.T) *auth.Login {
	db := testutil.DB(t)
	db.TruncateTables(t, testutil.TableUser)

	rdb := testutil.Redis(t)
	rdb.Flush(t)

	_ = fixtures.LoadUserFixtures(t, db.DB)

	repo := persistence.NewUserRepository(db.DB)
	hasher := security.NewPasswordHasher()
	jwt := security.NewJWTManager("TestSecretKey")
	refreshTokenRepo := persistence.NewRefreshTokenRepository(rdb.Client)
	refreshTokenManager := security.NewRefreshTokenManager()

	return auth.NewLogin(repo, hasher, jwt, refreshTokenRepo, refreshTokenManager)
}

func TestLogin_Success(t *testing.T) {
	login := setupLogin(t)

	input := auth.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "plainPassword",
	}

	response, err := login.Execute(context.Background(), input)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
}

func TestLogin_InvalidPassword(t *testing.T) {
	login := setupLogin(t)

	input := auth.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "wrongpassword",
	}

	response, err := login.Execute(context.Background(), input)
	assert.Error(t, err)
	assert.Empty(t, response.AccessToken)
	assert.Empty(t, response.RefreshToken)
}

func TestLogin_UserNotFound(t *testing.T) {
	login := setupLogin(t)

	input := auth.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password",
	}

	response, err := login.Execute(context.Background(), input)
	assert.Error(t, err)
	assert.Empty(t, response.AccessToken)
	assert.Empty(t, response.RefreshToken)
}
