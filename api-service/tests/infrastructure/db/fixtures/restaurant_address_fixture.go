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

func CreateRestaurantAddress(db *gorm.DB) (*restaurant.RestaurantAddress, error) {
	restAddr := restaurant.RestaurantAddress{
		House:      "12",
		Street:     "Danziger Str.",
		PostalCode: "10435",
		City:       "Hamburg",
		FullText:   "Danziger Str. 12, 10435 Berlin",
		Lat:        53.594970,
		Lon:        9.999253,
	}
	if err := db.Create(&restAddr).Error; err != nil {
		return nil, err
	}

	return &restAddr, nil
}
