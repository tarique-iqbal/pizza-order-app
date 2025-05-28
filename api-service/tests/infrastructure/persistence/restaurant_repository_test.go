package persistence_test

import (
	"api-service/internal/domain/restaurant"
	"api-service/internal/domain/user"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/infrastructure/db"
	"api-service/tests/infrastructure/db/fixtures"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type restaurantRepoTestEnv struct {
	User              *user.User
	RestaurantAddress *restaurant.RestaurantAddress
	RestaurantRepo    restaurant.RestaurantRepository
}

func setupRestaurantRepoTestEnv() restaurantRepoTestEnv {
	testDB := db.SetupTestDB()

	usr, err := fixtures.CreateUser(testDB, "Owner")
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

	return restaurantRepoTestEnv{
		User:              usr,
		RestaurantAddress: restAddr,
		RestaurantRepo:    restaurantRepo,
	}
}

func TestRestaurantRepository_Create(t *testing.T) {
	env := setupRestaurantRepoTestEnv()

	rest := restaurant.Restaurant{
		UserID:       env.User.ID,
		Name:         "Sushi Zen",
		Slug:         "sushi-zen",
		Email:        "info@sushizen.de",
		Phone:        "+49 89 99887766",
		AddressID:    env.RestaurantAddress.ID,
		DeliveryType: "third_party",
		DeliveryKm:   5,
		Specialties:  "italian",
		CreatedAt:    time.Now(),
	}

	err := env.RestaurantRepo.Create(context.Background(), &rest)
	assert.NoError(t, err)
	assert.NotZero(t, rest.ID)
}

func TestRestaurantRepository_FindBySlug(t *testing.T) {
	env := setupRestaurantRepoTestEnv()

	rest, err := env.RestaurantRepo.FindBySlug(context.Background(), "pizza-paradise")
	assert.NoError(t, err)
	assert.Equal(t, "pizza-paradise", rest.Slug)
}

func TestRestaurantRepository_IsOwnedBy(t *testing.T) {
	env := setupRestaurantRepoTestEnv()

	rest := restaurant.Restaurant{
		UserID:       env.User.ID,
		Name:         "Burger Meister",
		Slug:         "burger-meister",
		Email:        "contact@burgermeister.de",
		Phone:        "+49 351 22334455",
		AddressID:    env.RestaurantAddress.ID,
		DeliveryType: "pick_up",
		DeliveryKm:   7,
		Specialties:  "american",
		CreatedAt:    time.Now(),
	}

	err := env.RestaurantRepo.Create(context.Background(), &rest)
	assert.NoError(t, err)

	isOwner, err := env.RestaurantRepo.IsOwnedBy(context.Background(), rest.ID, rest.UserID)
	assert.NoError(t, err)
	assert.True(t, isOwner, "User is expected to be the owner")

	isOwner, err = env.RestaurantRepo.IsOwnedBy(context.Background(), rest.ID, 777)
	assert.NoError(t, err)
	assert.False(t, isOwner, "User is not expected to be the owner")

	isOwner, err = env.RestaurantRepo.IsOwnedBy(context.Background(), 888, rest.UserID)
	assert.NoError(t, err)
	assert.False(t, isOwner, "Non-existent restaurant is expected to return false")
}

func TestRestaurantRepository_IsSlugExists(t *testing.T) {
	env := setupRestaurantRepoTestEnv()

	exists, err := env.RestaurantRepo.IsSlugExists(context.Background(), "pizza-paradise")
	assert.NoError(t, err)
	assert.True(t, exists, "Slug is expected to be exists")

	exists, err = env.RestaurantRepo.IsSlugExists(context.Background(), "pizza-random")
	assert.NoError(t, err)
	assert.False(t, exists, "Slug is not expected to be exists")
}

func TestRestaurantRepository_IsEmailExists(t *testing.T) {
	env := setupRestaurantRepoTestEnv()

	exists, err := env.RestaurantRepo.IsEmailExists(context.Background(), "kontakt@pizzaparadise.de")
	assert.NoError(t, err)
	assert.True(t, exists, "Restaurant email is expected to be exists")

	exists, err = env.RestaurantRepo.IsEmailExists(context.Background(), "random@example.de")
	assert.NoError(t, err)
	assert.False(t, exists, "Restaurant email is not expected to be exists")
}
