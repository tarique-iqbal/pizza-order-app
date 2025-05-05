package validator

import (
	"api-service/internal/domain/restaurant"
	"api-service/internal/domain/user"
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	userRepo       user.UserRepository
	restaurantRepo restaurant.RestaurantRepository
}

func NewCustomValidator(
	userRepo user.UserRepository,
	restaurantRepo restaurant.RestaurantRepository,
) *CustomValidator {
	return &CustomValidator{userRepo: userRepo, restaurantRepo: restaurantRepo}
}

func (cv *CustomValidator) UniqueEmail(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		log.Printf("invalid kind: %s", fl.Field().Kind())
		return false
	}

	email := fl.Field().String()
	exists, err := cv.userRepo.EmailExists(email)
	if err != nil {
		log.Printf("error checking if email is unique: %v", err)
		return false
	}
	return !exists
}

func (cv *CustomValidator) UniqueRestaurantSlug(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		log.Printf("invalid kind: %s", fl.Field().Kind())
		return false
	}

	slug := fl.Field().String()
	exists, err := cv.restaurantRepo.SlugExists(slug)
	if err != nil {
		log.Printf("error checking if slug is unique: %v", err)
		return false
	}
	return !exists
}
