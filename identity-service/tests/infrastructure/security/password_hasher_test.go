package security_test

import (
	"identity-service/internal/domain/auth"
	"identity-service/internal/infrastructure/security"
	"testing"

	"github.com/stretchr/testify/assert"
)

func InitHasher() auth.PasswordHasher {
	return security.NewPasswordHasher()
}

func TestBcryptHasher_Hash(t *testing.T) {
	hasher := InitHasher()
	password := "securepassword"

	hashedPassword, err := hasher.Hash(password)

	assert.NoError(t, err, "Hashing should not return an error")
	assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
}

func TestBcryptHasher_Compare(t *testing.T) {
	hasher := InitHasher()
	password := "securepassword"

	hashedPassword, err := hasher.Hash(password)

	assert.NoError(t, err)
	assert.True(t, hasher.Compare(hashedPassword, password), "Valid password should return true")
	assert.False(t, hasher.Compare(hashedPassword, "wrongpassword"), "Invalid password should return false")
}
