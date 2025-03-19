package security_test

import (
	"pizza-order-api/internal/infrastructure/security"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "securepassword"

	hashedPassword, err := security.HashPassword(password)

	assert.NoError(t, err, "Hashing should not return an error")
	assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
}

func TestComparePassword(t *testing.T) {
	password := "securepassword"

	hashedPassword, err := security.HashPassword(password)

	assert.NoError(t, err)
	assert.True(t, security.ComparePassword(hashedPassword, password), "Valid password should return true")
	assert.False(t, security.ComparePassword(hashedPassword, "wrongpassword"), "Invalid password should return false")
}
