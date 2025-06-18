package restaurant_test

import (
	aRestaurant "api-service/internal/application/restaurant"
	"api-service/internal/domain/restaurant"
	"api-service/internal/domain/user"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/infrastructure/db"
	"api-service/tests/infrastructure/db/fixtures"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type createRestaurantUseCaseTestEnv struct {
	User               *user.User
	RestaurantAddress  *restaurant.RestaurantAddress
	CreateRestaurantUC *aRestaurant.CreateRestaurantUseCase
}

func setupCreateRestaurantUseCase() createRestaurantUseCaseTestEnv {
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

	restaurantRepo := persistence.NewRestaurantRepository(testDB)
	restAddrRepo := persistence.NewRestaurantAddressRepository(testDB)
	createRestaurantUC := aRestaurant.NewCreateRestaurantUseCase(testDB, restaurantRepo, restAddrRepo)

	return createRestaurantUseCaseTestEnv{
		User:               usr,
		RestaurantAddress:  restAddr,
		CreateRestaurantUC: createRestaurantUC,
	}
}

func TestCreateRestaurant_Success(t *testing.T) {
	env := setupCreateRestaurantUseCase()

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
	env := setupCreateRestaurantUseCase()

	existingEmail := "kontakt@pizzaparadise.de"

	input := aRestaurant.RestaurantCreateDTO{
		UserID: env.User.ID, Name: "Another", Email: existingEmail, Phone: "123",
		House: "1", Street: "X", City: "Y", PostalCode: "12345", DeliveryType: "pick_up", DeliveryKm: 5,
	}

	_, err := env.CreateRestaurantUC.Execute(context.Background(), input)
	assert.ErrorIs(t, err, restaurant.ErrEmailAlreadyExists)
}

func TestCreateRestaurant_DuplicateSlug(t *testing.T) {
	env := setupCreateRestaurantUseCase()

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
