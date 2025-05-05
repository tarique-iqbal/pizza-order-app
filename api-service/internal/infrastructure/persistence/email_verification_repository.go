package persistence

import (
	"api-service/internal/domain/auth"
	"context"

	"gorm.io/gorm"
)

type EmailVerificationRepositoryImpl struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) auth.EmailVerificationRepository {
	return &EmailVerificationRepositoryImpl{db: db}
}

func (repo *EmailVerificationRepositoryImpl) FindByEmail(
	ctx context.Context,
	email string,
) (*auth.EmailVerification, error) {
	var ev auth.EmailVerification
	err := repo.db.Where("email = ?", email).First(&ev).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &ev, nil
}

func (repo *EmailVerificationRepositoryImpl) Create(
	ctx context.Context,
	ev *auth.EmailVerification,
) error {
	return repo.db.Create(ev).Error
}

func (repo *EmailVerificationRepositoryImpl) Updates(
	ctx context.Context,
	ev *auth.EmailVerification,
) error {
	return repo.db.Updates(ev).Error
}
