package restaurant_test

import (
	aRestaurant "api-service/internal/application/restaurant"
	"api-service/internal/domain/user"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/internal/infrastructure/db"
	"api-service/tests/internal/infrastructure/db/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

type createRestaurantUseCaseTestEnv struct {
	User               *user.User
	CreateRestaurantUC *aRestaurant.CreateRestaurantUseCase
}

func setupCreateRestaurantUseCase() createRestaurantUseCaseTestEnv {
	testDB := db.SetupTestDB()

	usr, err := fixtures.CreateUser(testDB, "Owner")
	if err != nil {
		panic(err)
	}

	if err := fixtures.LoadRestaurantFixtures(testDB, usr); err != nil {
		panic(err)
	}

	restaurantRepo := persistence.NewRestaurantRepository(testDB)
	createRestaurantUC := aRestaurant.NewCreateRestaurantUseCase(restaurantRepo)

	return createRestaurantUseCaseTestEnv{
		User:               usr,
		CreateRestaurantUC: createRestaurantUC,
	}
}

func TestRestaurantUseCase_Create(t *testing.T) {
	env := setupCreateRestaurantUseCase()

	input := aRestaurant.RestaurantCreateDTO{
		UserID:  env.User.ID,
		Name:    "Test Restaurant",
		Slug:    "test-restaurant",
		Address: "123 Test Street",
	}

	rest, err := env.CreateRestaurantUC.Execute(input)

	assert.NoError(t, err)
	assert.NotZero(t, rest.ID)
	assert.Equal(t, env.User.ID, rest.UserID)
	assert.Equal(t, input.Name, rest.Name)
	assert.Equal(t, input.Slug, rest.Slug)
	assert.Equal(t, input.Address, rest.Address)
}
