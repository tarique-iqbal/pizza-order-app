package security_test

import (
	"os"
	"testing"
	"time"

	"pizza-order-api/internal/infrastructure/security"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Setenv("JWT_SECRET", "@TestSecret!")

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGenerateJWT(t *testing.T) {
	userID := uint(1)
	tokenString, err := security.GenerateJWT(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims, err := security.ParseJWT(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestValidateToken_ValidToken(t *testing.T) {
	userID := uint(1)
	tokenString, _ := security.GenerateJWT(userID)

	_, err := security.ValidateToken(tokenString)
	assert.NoError(t, err)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	_, err := security.ValidateToken("invalid.token.here")
	assert.Error(t, err)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	secretKey := os.Getenv("JWT_SECRET")

	expirationTime := time.Now().Add(-(24 * time.Hour))

	claims := &security.Claims{
		UserID: uint(1),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))

	_, err := security.ValidateToken(tokenString)
	assert.Error(t, err)
}
