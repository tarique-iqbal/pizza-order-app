package validator

import (
	"pizza-order-api/internal/domain/user"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	userRepo user.UserRepository
}

func NewCustomValidator(userRepo user.UserRepository) *CustomValidator {
	return &CustomValidator{userRepo: userRepo}
}

func (cv *CustomValidator) UniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	existingUser, _ := cv.userRepo.FindByEmail(email)
	return existingUser == nil
}
