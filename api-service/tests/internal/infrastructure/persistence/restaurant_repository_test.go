package persistence_test

import (
	"api-service/internal/domain/restaurant"
	"api-service/internal/domain/user"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/internal/infrastructure/db"
	"api-service/tests/internal/infrastructure/db/fixtures"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type restaurantRepoTestEnv struct {
	User           *user.User
	RestaurantRepo restaurant.RestaurantRepository
}

func setupRestaurantRepoTestEnv() restaurantRepoTestEnv {
	testDB := db.SetupTestDB()

	usr, err := fixtures.CreateUser(testDB, "Owner")
	if err != nil {
		panic(err)
	}

	if err := fixtures.LoadRestaurantFixtures(testDB, usr); err != nil {
		panic(err)
	}

	restaurantRepo := persistence.NewRestaurantRepository(testDB)

	return restaurantRepoTestEnv{
		User:           usr,
		RestaurantRepo: restaurantRepo,
	}
}

func TestRestaurantRepository_Create(t *testing.T) {
	env := setupRestaurantRepoTestEnv()

	r := restaurant.Restaurant{
		UserID:    env.User.ID,
		Name:      "Test Bistro",
		Slug:      "test-bistro",
		Address:   "789 Maple Street, Burger Town",
		CreatedAt: time.Now(),
	}

	err := env.RestaurantRepo.Create(&r)
	assert.NoError(t, err)
	assert.NotZero(t, r.ID)
}

func TestRestaurantRepository_FindBySlug(t *testing.T) {
	env := setupRestaurantRepoTestEnv()

	r, err := env.RestaurantRepo.FindBySlug("pizza-paradise")
	assert.NoError(t, err)
	assert.Equal(t, "Pizza Paradise", r.Name)
}
