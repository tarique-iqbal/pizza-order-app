package auth_test

import (
	"context"
	"identity-service/internal/application/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	"identity-service/tests/infrastructure/db"
	"identity-service/tests/infrastructure/db/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupSignInUseCase() *auth.SignInUseCase {
	testDB := db.SetupTestDB()

	if err := fixtures.LoadUserFixtures(testDB); err != nil {
		panic(err)
	}

	repo := persistence.NewUserRepository(testDB)
	hasher := security.NewPasswordHasher()
	jwt := security.NewJWTService("TestSecretKey")

	return auth.NewSignInUseCase(repo, hasher, jwt)
}

func TestSignInUseCase_Success(t *testing.T) {
	signInUC := setupSignInUseCase()

	token, err := signInUC.Execute(context.Background(), "john.doe@example.com", "plainPassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestSignInUseCase_InvalidPassword(t *testing.T) {
	signInUC := setupSignInUseCase()

	token, err := signInUC.Execute(context.Background(), "john.doe@example.com", "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestSignInUseCase_UserNotFound(t *testing.T) {
	signInUC := setupSignInUseCase()

	token, err := signInUC.Execute(context.Background(), "notfound@example.com", "password")
	assert.Error(t, err)
	assert.Empty(t, token)
}
