package persistence

import (
	"errors"
	"pizza-order-api/internal/domain/user"

	"gorm.io/gorm"
)

type UserRepositoryImplement struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepositoryImplement{DB: db}
}

func (repo *UserRepositoryImplement) FindByEmail(email string) (*user.User, error) {
	var user user.User
	result := repo.DB.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (repo *UserRepositoryImplement) Create(u *user.User) error {
	return repo.DB.Create(u).Error
}
