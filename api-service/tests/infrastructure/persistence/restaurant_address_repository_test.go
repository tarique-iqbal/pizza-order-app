package persistence_test

import (
	"api-service/internal/domain/restaurant"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/infrastructure/db"
	"api-service/tests/infrastructure/db/fixtures"
	"context"
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
