package persistence

import (
	"context"

	"api-service/internal/domain/restaurant"

	"gorm.io/gorm"
)

type restaurantAddressRepositoryImpl struct {
	db *gorm.DB
}

func NewRestaurantAddressRepository(db *gorm.DB) restaurant.RestaurantAddressRepository {
	return &restaurantAddressRepositoryImpl{db: db}
}

func (r *restaurantAddressRepositoryImpl) WithTx(
	tx *gorm.DB,
) restaurant.RestaurantAddressRepository {
	return &restaurantAddressRepositoryImpl{db: tx}
}

func (r *restaurantAddressRepositoryImpl) Create(
	ctx context.Context,
	address *restaurant.RestaurantAddress,
) error {
	return r.db.WithContext(ctx).Create(address).Error
}
