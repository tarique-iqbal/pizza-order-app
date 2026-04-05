package persistence

import (
	"context"
	"identity-service/internal/domain/user"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &GormUserRepository{db: db}
}

func (repo *GormUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := repo.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (repo *GormUserRepository) Create(ctx context.Context, u *user.User) error {
	return repo.db.WithContext(ctx).Create(u).Error
}

func (repo *GormUserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := repo.db.Model(&user.User{}).
		Where("email = ?", email).
		Count(&count).Error
	return count > 0, err
}
