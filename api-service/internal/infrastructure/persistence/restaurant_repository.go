package persistence

import (
	"api-service/internal/domain/restaurant"
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
