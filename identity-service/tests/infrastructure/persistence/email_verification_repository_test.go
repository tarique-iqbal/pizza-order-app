package persistence_test

import (
	"context"
	"identity-service/internal/domain/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/tests/infrastructure/db"
	"identity-service/tests/infrastructure/db/fixtures"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupEmailVerificationRepo() auth.EmailVerificationRepository {
	testDB := db.SetupTestDB()

	if err := fixtures.LoadEmailVerificationFixtures(testDB); err != nil {
		panic(err)
	}

	return persistence.NewEmailVerificationRepository(testDB)
}

func TestEmailVerificationRepository_Create(t *testing.T) {
	emVerRepo := setupEmailVerificationRepo()

	emailVerification := auth.EmailVerification{
		Email:     "adam.dangelo@example.com",
		Code:      "467923",
		IsUsed:    false,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
	}

	err := emVerRepo.Create(context.Background(), &emailVerification)

	assert.Nil(t, err)
	assert.NotZero(t, emailVerification.ID)
}

func TestEmailVerificationRepository_Updates(t *testing.T) {
	emVerRepo := setupEmailVerificationRepo()

	existing, _ := emVerRepo.FindByEmail(context.Background(), "john.doe@example.com")
	existing.Code = "478326"
	existing.IsUsed = false
	existing.ExpiresAt = time.Now().Add(15 * time.Minute)

	err := emVerRepo.Updates(context.Background(), existing)

	assert.Nil(t, err)
}

func TestEmailVerificationRepository_FindByEmail(t *testing.T) {
	emVerRepo := setupEmailVerificationRepo()

	emVer, err := emVerRepo.FindByEmail(context.Background(), "john.doe@example.com")

	assert.NoError(t, err)
	assert.Equal(t, "john.doe@example.com", emVer.Email)
}
