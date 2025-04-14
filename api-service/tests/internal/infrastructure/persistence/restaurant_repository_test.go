package persistence_test

import (
	"api-service/internal/domain/restaurant"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/internal/infrastructure/db"
	"api-service/tests/internal/infrastructure/db/fixtures"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var restaurantRepo restaurant.RestaurantRepository

func setupRestaurantRepo() restaurant.RestaurantRepository {
	testDB := db.SetupTestDB()

	if err := fixtures.LoadRestaurantFixtures(testDB); err != nil {
		panic(err)
	}

	return persistence.NewRestaurantRepository(testDB)
}

func TestRestaurantRepository_Create(t *testing.T) {
	restaurantRepo = setupRestaurantRepo()

	r := restaurant.Restaurant{
		UserID:    3,
		Name:      "Test Bistro",
		Slug:      "test-bistro",
		Address:   "789 Maple Street, Burger Town",
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	err := restaurantRepo.Create(&r)
	assert.NoError(t, err)
	assert.NotZero(t, r.ID)
}

func TestRestaurantRepository_FindBySlug(t *testing.T) {
	restaurantRepo = setupRestaurantRepo()

	r, err := restaurantRepo.FindBySlug("pizza-paradise")
	assert.NoError(t, err)
	assert.Equal(t, "Pizza Paradise", r.Name)
}
