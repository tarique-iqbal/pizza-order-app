package restaurant_test

import (
	aRestaurant "api-service/internal/application/restaurant"
	"api-service/internal/domain/restaurant"
	"api-service/internal/infrastructure/persistence"
	sErrors "api-service/internal/shared/errors"
	"api-service/tests/internal/infrastructure/db"
	"api-service/tests/internal/infrastructure/db/fixtures"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createPizzaSizeUseCaseTestEnv struct {
	Restaurant         *restaurant.Restaurant
	CreateRestaurantUC *aRestaurant.CreatePizzaSizeUseCase
}

func setupCreatePizzaSizeUseCase() createPizzaSizeUseCaseTestEnv {
	testDB := db.SetupTestDB()

	rest, err := fixtures.CreateRestaurant(testDB)
	if err != nil {
		panic(err)
	}

	if err := fixtures.LoadPizzaSizeFixtures(testDB, rest); err != nil {
		panic(err)
	}

	pizzaRepo := persistence.NewPizzaSizeRepository(testDB)
	restaurantRepo := persistence.NewRestaurantRepository(testDB)
	createPizzaSizeUC := aRestaurant.NewCreatePizzaSizeUseCase(pizzaRepo, restaurantRepo)

	return createPizzaSizeUseCaseTestEnv{
		Restaurant:         rest,
		CreateRestaurantUC: createPizzaSizeUC,
	}
}

func TestCreatePizzaSizeUseCase_Execute_Success(t *testing.T) {
	env := setupCreatePizzaSizeUseCase()

	input := aRestaurant.PizzaSizeCreateDTO{
		Title: "Large",
		Size:  12,
	}
	response, err := env.CreateRestaurantUC.Execute(
		context.Background(),
		env.Restaurant.ID,
		env.Restaurant.UserID,
		input,
	)

	require.NoError(t, err)
	assert.NotZero(t, response.ID)
	assert.Equal(t, "Large", response.Title)
	assert.Equal(t, 12, response.Size)
	assert.NotEmpty(t, response.CreatedAt)
}

func TestCreatePizzaSizeUseCase_Execute_Forbidden(t *testing.T) {
	env := setupCreatePizzaSizeUseCase()

	input := aRestaurant.PizzaSizeCreateDTO{
		Title: "Medium",
		Size:  10,
	}
	_, err := env.CreateRestaurantUC.Execute(
		context.Background(),
		env.Restaurant.ID,
		999,
		input,
	)

	assert.ErrorIs(t, err, sErrors.ErrForbidden)
}

func TestCreatePizzaSizeUseCase_Execute_PizzaSize_Duplicate(t *testing.T) {
	env := setupCreatePizzaSizeUseCase()

	input := aRestaurant.PizzaSizeCreateDTO{
		Title: "Large",
		Size:  12,
	}
	_, err := env.CreateRestaurantUC.Execute(
		context.Background(),
		env.Restaurant.ID,
		env.Restaurant.UserID,
		input,
	)
	require.NoError(t, err)

	_, err = env.CreateRestaurantUC.Execute(
		context.Background(),
		env.Restaurant.ID,
		env.Restaurant.UserID,
		input,
	)
	assert.ErrorIs(t, err, restaurant.ErrPizzaSizeAlreadyExists)
}
