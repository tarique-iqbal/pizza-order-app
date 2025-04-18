package security_test

import (
	"api-service/internal/domain/auth"
	"api-service/internal/infrastructure/security"
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
	tokenString, err := jwtService.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims, err := jwtService.ParseToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestJWTService_ValidToken(t *testing.T) {
	jwtService := InitJWT()

	userID := uint(1)
	tokenString, _ := jwtService.GenerateToken(userID)

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
