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
			Code:      "135864",
			IsUsed:    false,
			ExpiresAt: time.Now().Add(15 * time.Minute),
			CreatedAt: time.Now(),
		},
		{
			Email:     "adam.dangelo@example.com",
			Code:      "476190",
			IsUsed:    false,
			ExpiresAt: time.Now().Add(15 * time.Minute),
			CreatedAt: time.Now(),
		},
		{
			Email:     "alice@example.com",
			Code:      "347578",
			IsUsed:    false,
			ExpiresAt: time.Now().Add(15 * time.Minute),
			CreatedAt: time.Now(),
		},
		{
			Email:     "already.used@example.com",
			Code:      "137468",
			IsUsed:    true,
			ExpiresAt: time.Now().Add(15 * time.Minute),
			CreatedAt: time.Now(),
		},
		{
			Email:     "expired@example.com",
			Code:      "743802",
			IsUsed:    false,
			ExpiresAt: time.Now().Add(-1 * time.Minute),
			CreatedAt: time.Now(),
		},
	}

	for _, u := range emailVerification {
		db.Create(&u)
	}

	return nil
}
