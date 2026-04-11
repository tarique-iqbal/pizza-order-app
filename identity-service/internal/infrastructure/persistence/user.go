package persistence

import (
	"context"
	"errors"
	"identity-service/internal/domain/user"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepo{db: db}
}

func (repo *userRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
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

func (repo *userRepo) Create(ctx context.Context, u *user.User) error {
	return repo.db.WithContext(ctx).Create(u).Error
}

func (repo *userRepo) EmailExists(email string) (bool, error) {
	var count int64
	err := repo.db.Model(&user.User{}).
		Where("email = ?", email).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepo) FindByID(ctx context.Context, id int) (*user.User, error) {
	var u user.User

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Take(&u).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
