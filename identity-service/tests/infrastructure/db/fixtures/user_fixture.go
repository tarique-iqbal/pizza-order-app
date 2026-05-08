package fixtures

import (
	"identity-service/internal/domain/user"
	"identity-service/internal/infrastructure/security"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func LoadUserFixtures(t *testing.T, db *gorm.DB) error {
	hasher := security.NewPasswordHasher()
	password, err := hasher.Hash("plainPassword")
	require.NoError(t, err)

	users := []user.User{
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  password,
			Role:      "customer",
			CreatedAt: time.Now().UTC(),
		},
		{
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing@example.com",
			Password:  password,
			Role:      "owner",
			CreatedAt: time.Now().UTC(),
		},
	}

	for _, u := range users {
		userID, err := uuid.NewV7()
		require.NoError(t, err)

		u.ID = userID
		err = db.Create(&u).Error
		require.NoError(t, err)
	}

	return nil
}
