package persistence_test

import (
	"os"
	"pizza-order-api/internal/domain/restaurant"
	"pizza-order-api/internal/infrastructure/persistence"
	"pizza-order-api/tests/internal/infrastructure/db"
	"pizza-order-api/tests/internal/infrastructure/db/fixtures"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var restaurantRepo restaurant.RestaurantRepository

func TestMain(m *testing.M) {
	testDB = db.SetupTestDB()

	if err := fixtures.LoadRestaurantFixtures(testDB); err != nil {
		panic(err)
	}

	restaurantRepo = persistence.NewRestaurantRepository(testDB)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestRestaurantRepository_Create(t *testing.T) {
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
	r, err := restaurantRepo.FindBySlug("pizza-paradise")
	assert.NoError(t, err)
	assert.Equal(t, "Pizza Paradise", r.Name)
}
