package validator

import (
	"pizza-order-api/internal/domain/restaurant"
	"pizza-order-api/internal/domain/user"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	userRepo       user.UserRepository
	restaurantRepo restaurant.RestaurantRepository
}

func NewCustomValidator(userRepo user.UserRepository, restaurantRepo restaurant.RestaurantRepository) *CustomValidator {
	return &CustomValidator{userRepo: userRepo, restaurantRepo: restaurantRepo}
}

func (cv *CustomValidator) UniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	existingUser, _ := cv.userRepo.FindByEmail(email)
	return existingUser == nil
}

func (cv *CustomValidator) UniqueRestaurantSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	existingSlug, _ := cv.restaurantRepo.FindBySlug(slug)
	return existingSlug == nil
}
