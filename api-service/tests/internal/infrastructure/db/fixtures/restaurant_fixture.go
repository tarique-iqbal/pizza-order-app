package fixtures

import (
	"api-service/internal/domain/restaurant"
	"time"

	"gorm.io/gorm"
)

func LoadRestaurantFixtures(db *gorm.DB) error {
	restaurants := []restaurant.Restaurant{
		{
			UserID:    1,
			Name:      "Pizza Paradise",
			Slug:      "pizza-paradise",
			Address:   "123 Main Street, Food City",
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
		{
			UserID:    2,
			Name:      "Italiano Express",
			Slug:      "italiano-express",
			Address:   "456 Olive Avenue, Pasta Town",
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}

	for _, r := range restaurants {
		if err := db.Create(&r).Error; err != nil {
			return err
		}
	}

	return nil
}

func CreateRestaurant(db *gorm.DB) (*restaurant.Restaurant, error) {
	user, err := CreateUser(db)
	if err != nil {
		return nil, err
	}

	restaurant := restaurant.Restaurant{
		UserID:    user.ID,
		Name:      "Pizza Tonio",
		Slug:      "pizza-tonio",
		Address:   "123 Main Street, Food City",
		CreatedAt: time.Now(),
	}
	if err := db.Create(&restaurant).Error; err != nil {
		return nil, err
	}

	return &restaurant, nil
}
