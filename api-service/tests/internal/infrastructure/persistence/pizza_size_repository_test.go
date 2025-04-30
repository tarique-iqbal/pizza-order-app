package persistence_test

import (
	"api-service/internal/domain/restaurant"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/internal/infrastructure/db"
	"api-service/tests/internal/infrastructure/db/fixtures"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type pizzaSizeRepoTestEnv struct {
	Restaurant    *restaurant.Restaurant
	PizzaSizeRepo restaurant.PizzaSizeRepository
}

func setupPizzaSizeRepoTestEnv() pizzaSizeRepoTestEnv {
	testDB := db.SetupTestDB()

	rest, err := fixtures.CreateRestaurant(testDB)
	if err != nil {
		panic(err)
	}

	if err := fixtures.LoadPizzaSizeFixtures(testDB, rest); err != nil {
		panic(err)
	}

	pizzaSizeRepo := persistence.NewPizzaSizeRepository(testDB)

	return pizzaSizeRepoTestEnv{
		Restaurant:    rest,
		PizzaSizeRepo: pizzaSizeRepo,
	}
}

func TestPizzaSizeRepository_Create(t *testing.T) {
	env := setupPizzaSizeRepoTestEnv()

	ps := restaurant.PizzaSize{
		RestaurantID: env.Restaurant.ID,
		Title:        "Large",
		Size:         38,
		CreatedAt:    time.Now(),
	}
	err := env.PizzaSizeRepo.Create(context.Background(), &ps)

	assert.Nil(t, err)
	assert.NotZero(t, ps.ID)
	assert.Equal(t, env.Restaurant.ID, ps.RestaurantID)
}
