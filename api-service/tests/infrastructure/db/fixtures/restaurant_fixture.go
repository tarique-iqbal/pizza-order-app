package fixtures

import (
	"api-service/internal/domain/restaurant"
	"api-service/internal/domain/user"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func LoadRestaurantFixtures(db *gorm.DB, usr *user.User, restAddr *restaurant.RestaurantAddress) error {
	restaurants := []restaurant.Restaurant{
		{
			RestaurantUUID: uuid.New(),
			UserID:         usr.ID,
			Name:           "Pizza Paradise",
			Slug:           "pizza-paradise",
			Email:          "kontakt@pizzaparadise.de",
			Phone:          "+49 89 98765432",
			AddressID:      restAddr.ID,
			DeliveryType:   "own_delivery",
			DeliveryKm:     5,
			Specialties:    "italian,wood_fired",
			CreatedAt:      time.Now(),
		},
		{
			RestaurantUUID: uuid.New(),
			UserID:         usr.ID,
			Name:           "Anatolische Küche",
			Slug:           "anatolische-kueche",
			Email:          "kontakt@anatolisch.de",
			Phone:          "+49 40 76543210",
			AddressID:      restAddr.ID,
			DeliveryType:   "pick_up",
			DeliveryKm:     7,
			Specialties:    "italian",
			CreatedAt:      time.Now(),
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
	user, err := CreateUser(db, "owner")
	if err != nil {
		return nil, err
	}

	restAddr, err := CreateRestaurantAddress(db)
	if err != nil {
		panic(err)
	}

	restaurant := restaurant.Restaurant{
		RestaurantUUID: uuid.New(),
		UserID:         user.ID,
		Name:           "Pizza Tonio",
		Slug:           "pizza-tonio",
		Email:          "hallo@pizzatonio.de",
		Phone:          "+49 69 22334455",
		AddressID:      restAddr.ID,
		DeliveryType:   "third_party",
		DeliveryKm:     6,
		Specialties:    "italian",
		CreatedAt:      time.Now(),
	}
	if err := db.Create(&restaurant).Error; err != nil {
		return nil, err
	}

	return &restaurant, nil
}
