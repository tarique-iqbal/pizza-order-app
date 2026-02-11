package security_test

import (
	"restaurant-service/internal/domain/auth"
	"restaurant-service/internal/infrastructure/security"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jwtClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func InitJWT() auth.JWTService {
	return security.NewJWTService("TestSecretKey")
}

func TestJWTService_GenerateToken(t *testing.T) {
	jwtService := InitJWT()

	userID := uint(1)
	role := "user"
	tokenString, err := jwtService.GenerateToken(userID, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims, err := jwtService.ParseToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
}

func TestJWTService_ValidToken(t *testing.T) {
	jwtService := InitJWT()

	userID := uint(1)
	role := "owner"
	tokenString, _ := jwtService.GenerateToken(userID, role)

	_, err := jwtService.ParseToken(tokenString)
	assert.NoError(t, err)
}

func TestJWTService_InvalidToken(t *testing.T) {
	jwtService := InitJWT()

	_, err := jwtService.ParseToken("invalid.token.here")
	assert.Error(t, err)
}

func TestJWTService_ExpiredToken(t *testing.T) {
	jwtService := InitJWT()

	expiredClaims := jwtClaims{
		UserID: 123,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, _ := token.SignedString([]byte("TestSecretKey"))

	_, err := jwtService.ParseToken(tokenString)

	require.Error(t, err)
	require.Contains(t, err.Error(), "token is expired")
}
