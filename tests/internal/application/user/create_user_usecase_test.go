package user_test

import (
	aUser "pizza-order-api/internal/application/user"
	"pizza-order-api/internal/infrastructure/persistence"
	"pizza-order-api/tests/internal/infrastructure/db"
	"pizza-order-api/tests/internal/infrastructure/db/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createUseCase *aUser.CreateUserUseCase

func createUserUseCase() *aUser.CreateUserUseCase {
	testDB := db.SetupTestDB()

	if err := fixtures.LoadUserFixtures(testDB); err != nil {
		panic(err)
	}

	userRepo := persistence.NewUserRepository(testDB)
	return aUser.NewCreateUserUseCase(userRepo)
}

func TestCreateUserUseCase(t *testing.T) {
	createUseCase = createUserUseCase()

	input := aUser.UserCreateDTO{
		FirstName: "Adam",
		LastName:  "D'Angelo",
		Email:     "adam.dangelo@example.com",
		Password:  "securepassword",
	}

	user, err := createUseCase.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Adam", user.FirstName)
}
