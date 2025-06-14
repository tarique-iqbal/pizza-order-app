package persistence_test

import (
	"api-service/internal/domain/user"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/infrastructure/db"
	"api-service/tests/infrastructure/db/fixtures"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupUserRepo() user.UserRepository {
	testDB := db.SetupTestDB()

	if err := fixtures.LoadUserFixtures(testDB); err != nil {
		panic(err)
	}

	return persistence.NewUserRepository(testDB)
}

func TestUserRepository_Create(t *testing.T) {
	userRepo := setupUserRepo()

	usr := user.User{
		FirstName: "Adam",
		LastName:  "D'Angelo",
		Email:     "adam.dangelo@example.com",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now(),
	}

	err := userRepo.Create(context.Background(), &usr)

	assert.Nil(t, err)
	assert.NotZero(t, usr.ID)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	userRepo := setupUserRepo()

	usr, err := userRepo.FindByEmail(context.Background(), "john.doe@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "John", usr.FirstName)
}

func TestUserRepository_EmailExists(t *testing.T) {
	userRepo := setupUserRepo()

	exists, err := userRepo.EmailExists("john.doe@example.com")
	assert.NoError(t, err)
	assert.True(t, exists, "Email is expected to be exists")

	exists, err = userRepo.EmailExists("random@example.com")
	assert.NoError(t, err)
	assert.False(t, exists, "Email is not expected to be exists")
}
