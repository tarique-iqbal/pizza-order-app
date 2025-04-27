package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	dAuth "api-service/internal/domain/auth"
	iAuth "api-service/internal/infrastructure/auth"
	"api-service/internal/infrastructure/persistence"
	"api-service/tests/internal/infrastructure/db"
	"api-service/tests/internal/infrastructure/db/fixtures"
)

func setupCodeVerificationService() dAuth.CodeVerifier {
	testDB := db.SetupTestDB()
	repo := persistence.NewEmailVerificationRepository(testDB)

	if err := fixtures.LoadEmailVerificationFixtures(testDB); err != nil {
		panic(err)
	}

	return iAuth.NewCodeVerificationService(repo)
}

func TestVerify_Success(t *testing.T) {
	svc := setupCodeVerificationService()

	err := svc.Verify("alice@example.com", "347578")
	assert.NoError(t, err)
}

func TestVerify_CodeMismatch(t *testing.T) {
	svc := setupCodeVerificationService()

	err := svc.Verify("alice@example.com", "010101")
	assert.ErrorIs(t, err, dAuth.ErrCodeInvalid)
}

func TestVerify_AlreadyUsed(t *testing.T) {
	svc := setupCodeVerificationService()

	err := svc.Verify("already.used@example.com", "137468")
	assert.ErrorIs(t, err, dAuth.ErrCodeUsed)
}

func TestVerify_Expired(t *testing.T) {
	svc := setupCodeVerificationService()

	err := svc.Verify("expired@example.com", "743802")
	assert.ErrorIs(t, err, dAuth.ErrCodeExpired)
}

func TestVerify_CodeNotIssued(t *testing.T) {
	svc := setupCodeVerificationService()

	err := svc.Verify("not.found@example.com", "578578")
	assert.ErrorIs(t, err, dAuth.ErrCodeNotIssued)
}
