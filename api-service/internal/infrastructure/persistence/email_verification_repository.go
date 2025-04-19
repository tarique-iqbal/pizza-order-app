package persistence

import (
	"api-service/internal/domain/auth"

	"gorm.io/gorm"
)

type EmailVerificationRepositoryImpl struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) auth.EmailVerificationRepository {
	return &EmailVerificationRepositoryImpl{db: db}
}

func (repo *EmailVerificationRepositoryImpl) FindByEmail(email string) (*auth.EmailVerification, error) {
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

func (repo *EmailVerificationRepositoryImpl) Create(ev *auth.EmailVerification) error {
	return repo.db.Create(ev).Error
}

func (repo *EmailVerificationRepositoryImpl) Updates(ev *auth.EmailVerification) error {
	return repo.db.Updates(ev).Error
}
