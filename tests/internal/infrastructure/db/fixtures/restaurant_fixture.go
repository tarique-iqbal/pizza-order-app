package fixtures

import (
	"pizza-order-api/internal/domain/restaurant"
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
