package fixtures

import (
	"api-service/internal/domain/auth"
	"time"

	"gorm.io/gorm"
)

func LoadEmailVerificationFixtures(db *gorm.DB) error {
	emailVerification := []auth.EmailVerification{
		{
			Email:     "john.doe@example.com",
			Code:      "036934",
			IsUsed:    false,
			ExpiresAt: time.Now().Add(15 * time.Minute),
			CreatedAt: time.Now(),
		},
	}

	for _, u := range emailVerification {
		db.Create(&u)
	}

	return nil
}
