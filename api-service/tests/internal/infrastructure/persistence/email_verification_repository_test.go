package persistence_test

import (
	"api-service/internal/domain/auth"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/internal/infrastructure/db"
	"api-service/tests/internal/infrastructure/db/fixtures"
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
	var repo auth.EmailVerificationRepository = setupEmailVerificationRepo()

	emailVerification := &auth.EmailVerification{
		Email:     "adam.dangelo@example.com",
		Code:      "467923",
		IsUsed:    false,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
	}

	err := repo.Create(emailVerification)

	assert.Nil(t, err)
	assert.NotZero(t, emailVerification.ID)
}

func TestEmailVerificationRepository_Updates(t *testing.T) {
	var repo auth.EmailVerificationRepository = setupEmailVerificationRepo()

	existing, _ := repo.FindByEmail("john.doe@example.com")
	existing.Code = "478326"
	existing.IsUsed = false
	existing.ExpiresAt = time.Now().Add(15 * time.Minute)

	err := repo.Updates(existing)

	assert.Nil(t, err)
}

func TestEmailVerificationRepository_FindByEmail(t *testing.T) {
	var repo auth.EmailVerificationRepository = setupEmailVerificationRepo()

	r, err := repo.FindByEmail("john.doe@example.com")

	assert.NoError(t, err)
	assert.Equal(t, "john.doe@example.com", r.Email)
}
