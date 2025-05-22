package fixtures

import (
	"api-service/internal/domain/restaurant"

	"gorm.io/gorm"
)

func LoadRestaurantAddressFixtures(db *gorm.DB) error {
	restAddrs := []restaurant.RestaurantAddress{
		{
			House:      "17",
			Street:     "Winterhuder Marktplatz",
			PostalCode: "22299",
			City:       "Hamburg",
			FullText:   "Winterhuder Marktplatz 17, 22299 Hamburg",
			Lat:        53.594970,
			Lon:        9.999253,
		},
	}

	for _, u := range restAddrs {
		db.Create(&u)
	}

	return nil
}
