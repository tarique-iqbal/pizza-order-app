package persistence

import (
	"api-service/internal/domain/restaurant"
	"context"
	"errors"

	"gorm.io/gorm"
)

type RestaurantRepositoryImpl struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) restaurant.RestaurantRepository {
	return &RestaurantRepositoryImpl{db: db}
}

func (r *RestaurantRepositoryImpl) Create(res *restaurant.Restaurant) error {
	return r.db.Create(res).Error
}

func (repo *RestaurantRepositoryImpl) FindBySlug(slug string) (*restaurant.Restaurant, error) {
	var r restaurant.Restaurant
	err := repo.db.Where("slug = ?", slug).First(&r).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &r, err
}

func (repo *RestaurantRepositoryImpl) IsOwnedBy(ctx context.Context, restaurantID uint, ownerID uint) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&restaurant.Restaurant{}).
		Where("id = ? AND user_id = ?", restaurantID, ownerID).
		Count(&count).Error
	return count > 0, err
}
