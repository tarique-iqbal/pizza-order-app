package restaurant_test

import (
	"context"
	"fmt"
	aRestaurant "restaurant-service/internal/application/restaurant"
	"restaurant-service/internal/domain/restaurant"
	"restaurant-service/internal/domain/user"
	"restaurant-service/internal/infrastructure/persistence"
	"restaurant-service/tests/infrastructure/db"
	"restaurant-service/tests/infrastructure/db/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockGeocoder struct {
	lat float64
	lon float64
	err error
}

func (m *mockGeocoder) GeocodeAddress(addr restaurant.RestaurantAddress) (float64, float64, error) {
	return m.lat, m.lon, m.err
}

type createRestaurantUseCaseTestEnv struct {
	User               *user.User
	RestaurantAddress  *restaurant.RestaurantAddress
	CreateRestaurantUC *aRestaurant.CreateRestaurantUseCase
}

func setupCreateRestaurantUseCase(lat float64, lon float64, errGeo error) createRestaurantUseCaseTestEnv {
	testDB := db.SetupTestDB()

	usr, err := fixtures.CreateUser(testDB, "owner")
	if err != nil {
		panic(err)
	}

	restAddr, err := fixtures.CreateRestaurantAddress(testDB)
	if err != nil {
		panic(err)
	}

	if err := fixtures.LoadRestaurantFixtures(testDB, usr, restAddr); err != nil {
		panic(err)
	}

	mockGeo := &mockGeocoder{lat: lat, lon: lon, err: errGeo}
	restaurantRepo := persistence.NewRestaurantRepository(testDB)
	restAddrRepo := persistence.NewRestaurantAddressRepository(testDB)
	createRestaurantUC := aRestaurant.NewCreateRestaurantUseCase(testDB, mockGeo, restaurantRepo, restAddrRepo)

	return createRestaurantUseCaseTestEnv{
		User:               usr,
		RestaurantAddress:  restAddr,
		CreateRestaurantUC: createRestaurantUC,
	}
}

func TestCreateRestaurant_Success(t *testing.T) {
	env := setupCreateRestaurantUseCase(52.52, 13.405, nil)

	input := aRestaurant.RestaurantCreateDTO{
		UserID:       env.User.ID,
		Name:         "Test Restaurant",
		Email:        "unique@test.com",
		Phone:        "+49 89 22334455",
		House:        "1",
		Street:       "Main Str.",
		City:         "Cityville",
		PostalCode:   "12345",
		DeliveryType: "own_delivery",
		DeliveryKm:   5,
		Specialties:  []string{"italian", "wood_fired"},
	}

	rest, err := env.CreateRestaurantUC.Execute(context.Background(), input)
	assert.NoError(t, err)
	assert.NotZero(t, rest.ID)
	assert.Equal(t, input.Name, rest.Name)
}

func TestCreateRestaurant_DuplicateEmail(t *testing.T) {
	env := setupCreateRestaurantUseCase(52.52, 13.405, nil)

	existingEmail := "kontakt@pizzaparadise.de"

	input := aRestaurant.RestaurantCreateDTO{
		UserID: env.User.ID, Name: "Another", Email: existingEmail, Phone: "123",
		House: "1", Street: "X", City: "Y", PostalCode: "12345", DeliveryType: "pick_up", DeliveryKm: 5,
	}

	_, err := env.CreateRestaurantUC.Execute(context.Background(), input)
	assert.ErrorIs(t, err, restaurant.ErrEmailAlreadyExists)
}

func TestCreateRestaurant_DuplicateSlug(t *testing.T) {
	env := setupCreateRestaurantUseCase(52.52, 13.405, nil)

	name := "Pizza Paradise"
	city := "Hamburg"

	input := aRestaurant.RestaurantCreateDTO{
		UserID: env.User.ID, Name: name, Email: "newemail@test.com", Phone: "123",
		House: "1", Street: "X", City: city, PostalCode: "12345", DeliveryType: "pick_up", DeliveryKm: 5,
	}

	rest, err := env.CreateRestaurantUC.Execute(context.Background(), input)
	assert.NoError(t, err)
	assert.NotEqual(t, fmt.Sprintf("%s-%s", name, city), rest.Slug) // should be unique
}
