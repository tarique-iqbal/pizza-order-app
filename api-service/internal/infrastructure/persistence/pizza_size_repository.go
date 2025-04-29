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
