package fixtures

import (
	"api-service/internal/domain/restaurant"
	"time"

	"gorm.io/gorm"
)

func LoadPizzaSizeFixtures(db *gorm.DB, rest *restaurant.Restaurant) error {
	pizzaSizes := []restaurant.PizzaSize{
		{
			RestaurantID: rest.ID,
			Title:        "Classic",
			Size:         26,
			CreatedAt:    time.Now(),
		},
		{
			RestaurantID: rest.ID,
			Title:        "Medium",
			Size:         32,
			CreatedAt:    time.Now(),
		},
	}

	for _, ps := range pizzaSizes {
		if err := db.Create(&ps).Error; err != nil {
			return err
		}
	}

	return nil
}
