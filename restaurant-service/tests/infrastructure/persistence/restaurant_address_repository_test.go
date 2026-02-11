package persistence_test

import (
	"context"
	"restaurant-service/internal/domain/restaurant"
	"restaurant-service/internal/infrastructure/persistence"
	"restaurant-service/tests/infrastructure/db"
	"restaurant-service/tests/infrastructure/db/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupRestaurantAddressRepo() restaurant.RestaurantAddressRepository {
	testDB := db.SetupTestDB()

	if err := fixtures.LoadRestaurantAddressFixtures(testDB); err != nil {
		panic(err)
	}

	return persistence.NewRestaurantAddressRepository(testDB)
}

func TestRestaurantAddressRepository_Create(t *testing.T) {
	restaurantAddressRepo := setupRestaurantAddressRepo()

	restAddr := restaurant.RestaurantAddress{
		House:      "77",
		Street:     "Langenfelder Damm",
		PostalCode: "22525",
		City:       "Hamburg",
		FullText:   "Langenfelder Damm 77, 22525 Hamburg",
		Lat:        53.581692,
		Lon:        9.936925,
	}

	err := restaurantAddressRepo.Create(context.Background(), &restAddr)

	assert.Nil(t, err)
	assert.NotZero(t, restAddr.ID)
}
