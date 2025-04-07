package persistence

import (
	"pizza-order-api/internal/domain/user"

	"gorm.io/gorm"
)

type UserRepositoryImplement struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepositoryImplement{db: db}
}

func (repo *UserRepositoryImplement) FindByEmail(email string) (*user.User, error) {
	var u user.User
	err := repo.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (repo *UserRepositoryImplement) Create(u *user.User) error {
	return repo.db.Create(u).Error
}
