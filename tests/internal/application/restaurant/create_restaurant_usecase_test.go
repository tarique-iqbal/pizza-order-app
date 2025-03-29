package restaurant_test

import (
	"os"
	aRestaurant "pizza-order-api/internal/application/restaurant"
	"pizza-order-api/internal/infrastructure/persistence"
	"pizza-order-api/tests/internal/infrastructure/db"
	"pizza-order-api/tests/internal/infrastructure/db/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var createUseCase *aRestaurant.CreateRestaurantUseCase

func TestMain(m *testing.M) {
	testDB = db.SetupTestDB()

	if err := fixtures.LoadRestaurantFixtures(testDB); err != nil {
		panic(err)
	}

	restaurantRepo := persistence.NewRestaurantRepository(testDB)
	createUseCase = aRestaurant.NewCreateRestaurantUseCase(restaurantRepo)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestRestaurantUseCase_Create(t *testing.T) {
	input := aRestaurant.RestaurantCreateDTO{
		UserID:  1,
		Name:    "Test Restaurant",
		Slug:    "test-restaurant",
		Address: "123 Test Street",
	}

	createdRestaurant, err := createUseCase.Execute(input)

	assert.NoError(t, err)
	assert.NotZero(t, createdRestaurant.ID)
	assert.Equal(t, input.Name, createdRestaurant.Name)
	assert.Equal(t, input.Slug, createdRestaurant.Slug)
	assert.Equal(t, input.Address, createdRestaurant.Address)
}
