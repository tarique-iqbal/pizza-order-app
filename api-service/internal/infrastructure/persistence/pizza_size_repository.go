package persistence

import (
	"context"

	"api-service/internal/domain/restaurant"

	"gorm.io/gorm"
)

type pizzaSizeRepositoryImpl struct {
	db *gorm.DB
}

func NewPizzaSizeRepository(db *gorm.DB) restaurant.PizzaSizeRepository {
	return &pizzaSizeRepositoryImpl{db: db}
}

func (repo *pizzaSizeRepositoryImpl) Create(ctx context.Context, size *restaurant.PizzaSize) error {
	return repo.db.WithContext(ctx).Create(size).Error
}

func (repo *pizzaSizeRepositoryImpl) PizzaSizeExists(
	ctx context.Context,
	restaurantID uint,
	size int,
) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&restaurant.PizzaSize{}).
		Where("restaurant_id = ? AND size = ?", restaurantID, size).
		Count(&count).Error
	return count > 0, err
}
