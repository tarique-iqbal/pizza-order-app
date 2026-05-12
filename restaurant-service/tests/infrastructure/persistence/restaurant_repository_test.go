package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"restaurant-service/internal/domain/restaurant"
	"restaurant-service/internal/infrastructure/persistence"
	"restaurant-service/tests/infrastructure/db/fixtures"
	"restaurant-service/tests/testutil"
)

type restaurantRepoEnv struct {
	DB             *gorm.DB
	RestaurantRepo restaurant.RestaurantRepository
}

func setupRestaurantRepoEnv(t *testing.T) restaurantRepoEnv {
	db := testutil.DB(t)
	db.TruncateTables(t, testutil.TableRestaurant)

	_ = fixtures.LoadRestaurantFixtures(t, db.DB)

	restaurantRepo := persistence.NewRestaurantRepository(db.DB)

	return restaurantRepoEnv{
		DB:             db.DB,
		RestaurantRepo: restaurantRepo,
	}
}

func TestRestaurantRepository_Create(t *testing.T) {
	env := setupRestaurantRepoEnv(t)

	checklist := restaurant.NewChecklist()
	checklist.Complete(restaurant.ChecklistBasic)

	res := restaurant.Restaurant{
		ID:        testutil.MustNewID(),
		OwnerID:   testutil.MustNewID(),
		Name:      "Pizza Paradise",
		VATNumber: "DE323678654",
		Checklist: checklist,
		CreatedAt: time.Now().UTC(),
	}

	err := env.RestaurantRepo.Create(context.Background(), &res)
	assert.NoError(t, err)
	assert.NotZero(t, res.ID)
	assert.True(t, res.Checklist[restaurant.ChecklistBasic])
}

func TestRestaurantRepository_FindBySlug(t *testing.T) {
	env := setupRestaurantRepoEnv(t)

	res, err := env.RestaurantRepo.FindBySlug(context.Background(), "anatolische-kueche")
	assert.NoError(t, err)
	assert.Equal(t, "anatolische-kueche", *res.Slug)
}

func TestRestaurantRepository_IsOwnedBy(t *testing.T) {
	env := setupRestaurantRepoEnv(t)

	var res restaurant.Restaurant
	require.NoError(t, env.DB.Last(&res).Error)

	isOwner, err := env.RestaurantRepo.IsOwnedBy(context.Background(), res.ID, res.OwnerID)
	assert.NoError(t, err)
	assert.True(t, isOwner, "User is expected to be the owner")

	isOwner, err = env.RestaurantRepo.IsOwnedBy(context.Background(), res.ID, testutil.MustNewID())
	assert.NoError(t, err)
	assert.False(t, isOwner, "User is not expected to be the owner")

	isOwner, err = env.RestaurantRepo.IsOwnedBy(context.Background(), testutil.MustNewID(), res.OwnerID)
	assert.NoError(t, err)
	assert.False(t, isOwner, "Non-existent restaurant is expected to return false")
}

func TestRestaurantRepository_IsSlugExists(t *testing.T) {
	env := setupRestaurantRepoEnv(t)

	exists, err := env.RestaurantRepo.IsSlugExists(context.Background(), "anatolische-kueche")
	assert.NoError(t, err)
	assert.True(t, exists, "Slug is expected to be exists")

	exists, err = env.RestaurantRepo.IsSlugExists(context.Background(), "pizza-random")
	assert.NoError(t, err)
	assert.False(t, exists, "Slug is not expected to be exists")
}

func TestRestaurantRepository_IsEmailExists(t *testing.T) {
	env := setupRestaurantRepoEnv(t)

	exists, err := env.RestaurantRepo.IsEmailExists(context.Background(), "kontakt@anatolisch.de")
	assert.NoError(t, err)
	assert.True(t, exists, "Restaurant email is expected to be exists")

	exists, err = env.RestaurantRepo.IsEmailExists(context.Background(), "random@example.de")
	assert.NoError(t, err)
	assert.False(t, exists, "Restaurant email is not expected to be exists")
}
