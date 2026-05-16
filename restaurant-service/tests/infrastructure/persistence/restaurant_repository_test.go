package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestRestaurantRepository_Update(t *testing.T) {
	env := setupRestaurantRepoEnv(t)

	var res restaurant.Restaurant
	err := env.DB.First(&res).Error
	assert.NoError(t, err)

	deliveryKm := int16(5)

	res.DeliveryKm = &deliveryKm
	res.DeliveryType = restaurant.DeliveryOwn
	res.Checklist.Complete(restaurant.ChecklistDelivery)

	err = env.RestaurantRepo.Update(context.Background(), &res)
	assert.NoError(t, err)

	var r restaurant.Restaurant
	err = env.DB.Take(&r, "id = ?", res.ID).Error
	assert.NoError(t, err)

	assert.NotNil(t, r.DeliveryKm)
	assert.Equal(t, int16(5), *r.DeliveryKm)
	assert.Equal(t, restaurant.DeliveryOwn, r.DeliveryType)
	assert.True(t, r.Checklist[restaurant.ChecklistDelivery])
}

func TestRestaurantRepository_FindBySlug(t *testing.T) {
	env := setupRestaurantRepoEnv(t)

	res, err := env.RestaurantRepo.FindBySlug(
		context.Background(),
		"anatolische-kueche", // from fixture
	)
	assert.NoError(t, err)

	assert.NotNil(t, res)
	assert.NotNil(t, res.Slug)
	assert.Equal(t, "anatolische-kueche", *res.Slug)

	res, err = env.RestaurantRepo.FindBySlug(
		context.Background(),
		"not-exist",
	)
	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestRestaurantRepository_FindByIDAndOwner(t *testing.T) {
	env := setupRestaurantRepoEnv(t)

	var existing restaurant.Restaurant
	err := env.DB.First(&existing).Error
	assert.NoError(t, err)

	tests := []struct {
		name         string
		restaurantID uuid.UUID
		ownerID      uuid.UUID
		found        bool
	}{
		{
			name:         "found",
			restaurantID: existing.ID,
			ownerID:      existing.OwnerID,
			found:        true,
		},
		{
			name:         "wrong owner",
			restaurantID: existing.ID,
			ownerID:      testutil.MustNewID(),
		},
		{
			name:         "wrong restaurant id",
			restaurantID: testutil.MustNewID(),
			ownerID:      existing.OwnerID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := env.RestaurantRepo.FindByIDAndOwner(
				context.Background(),
				tt.restaurantID,
				tt.ownerID,
			)

			assert.NoError(t, err)

			if !tt.found {
				assert.Nil(t, res)
				return
			}

			assert.NotNil(t, res)
			assert.Equal(t, existing.ID, res.ID)
			assert.Equal(t, existing.OwnerID, res.OwnerID)
		})
	}
}

func TestRestaurantRepository_IsOwner(t *testing.T) {
	env := setupRestaurantRepoEnv(t)

	var res restaurant.Restaurant
	require.NoError(t, env.DB.Last(&res).Error)

	isOwner, err := env.RestaurantRepo.IsOwner(context.Background(), res.ID, res.OwnerID)
	assert.NoError(t, err)
	assert.True(t, isOwner, "User is expected to be the owner")

	isOwner, err = env.RestaurantRepo.IsOwner(context.Background(), res.ID, testutil.MustNewID())
	assert.NoError(t, err)
	assert.False(t, isOwner, "User is not expected to be the owner")

	isOwner, err = env.RestaurantRepo.IsOwner(context.Background(), testutil.MustNewID(), res.OwnerID)
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
