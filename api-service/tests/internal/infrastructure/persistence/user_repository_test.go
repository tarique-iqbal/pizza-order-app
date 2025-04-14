package persistence_test

import (
	"api-service/internal/domain/user"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/internal/infrastructure/db"
	"api-service/tests/internal/infrastructure/db/fixtures"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var userRepo user.UserRepository

func setupUserRepo() user.UserRepository {
	testDB := db.SetupTestDB()

	if err := fixtures.LoadUserFixtures(testDB); err != nil {
		panic(err)
	}

	return persistence.NewUserRepository(testDB)
}

func TestUserRepository_Create(t *testing.T) {
	userRepo = setupUserRepo()

	newUser := &user.User{
		FirstName: "Adam",
		LastName:  "D'Angelo",
		Email:     "adam.dangelo@example.com",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	err := userRepo.Create(newUser)

	assert.Nil(t, err)
	assert.NotZero(t, newUser.ID)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	userRepo = setupUserRepo()

	r, err := userRepo.FindByEmail("john.doe@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "John", r.FirstName)
}
