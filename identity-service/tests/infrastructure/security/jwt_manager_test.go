package security_test

import (
	"identity-service/internal/domain/auth"
	"identity-service/internal/infrastructure/security"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jwtClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func InitJWT() auth.JWTManager {
	return security.NewJWTManager("TestSecretKey")
}

func TestJWTManager_GenerateToken(t *testing.T) {
	jwtManager := InitJWT()

	userID, _ := uuid.NewV7()
	role := "customer"

	tokenString, err := jwtManager.Generate(userID.String(), role)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims, err := jwtManager.Parse(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), claims.UserID)
	assert.Equal(t, role, claims.Role)
}

func TestJWTManager_ValidToken(t *testing.T) {
	jwtManager := InitJWT()

	userID, _ := uuid.NewV7()
	role := "owner"

	tokenString, _ := jwtManager.Generate(userID.String(), role)

	_, err := jwtManager.Parse(tokenString)
	assert.NoError(t, err)
}

func TestJWTManager_InvalidToken(t *testing.T) {
	jwtManager := InitJWT()

	_, err := jwtManager.Parse("invalid.token.here")
	assert.Error(t, err)
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	jwtManager := InitJWT()

	userID, _ := uuid.NewV7()

	expiredClaims := jwtClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, _ := token.SignedString([]byte("TestSecretKey"))

	_, err := jwtManager.Parse(tokenString)

	require.Error(t, err)
	require.Contains(t, err.Error(), "token is expired")
}
