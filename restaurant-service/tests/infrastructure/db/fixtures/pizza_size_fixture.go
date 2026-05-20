package fixtures

import (
	"time"

	"gorm.io/gorm"

	"restaurant-service/internal/domain/restaurant"
)

func LoadPizzaSizeFixtures(db *gorm.DB, rest *restaurant.Restaurant) error {
	pizzaSizes := []restaurant.PizzaSize{
		{
			DiameterCm: 22,
			CreatedAt:  time.Now().UTC(),
		},
		{
			DiameterCm: 26,
			CreatedAt:  time.Now().UTC(),
		},
		{
			DiameterCm: 30,
			CreatedAt:  time.Now().UTC(),
		},
		{
			DiameterCm: 34,
			CreatedAt:  time.Now().UTC(),
		},
		{
			DiameterCm: 36,
			CreatedAt:  time.Now().UTC(),
		},
	}

	for _, ps := range pizzaSizes {
		if err := db.Create(&ps).Error; err != nil {
			return err
		}
	}

	return nil
}
