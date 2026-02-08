package middlewares_test

import (
	"identity-service/internal/domain/auth"
	"identity-service/internal/infrastructure/security"
	"identity-service/internal/interfaces/http/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func InitJWT() auth.JWTService {
	gin.SetMode(gin.TestMode)
	return security.NewJWTService("TestSecretKey")
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtService := InitJWT()

	userID := uint(1)
	role := "user"
	token, _ := jwtService.GenerateToken(userID, role)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	ctx.Request = req

	middlewares.AuthMiddleware(jwtService)(ctx)
	ctxUserID := ctx.MustGet("userID").(uint)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, userID, ctxUserID)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwtService := InitJWT()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	ctx.Request = req

	middlewares.AuthMiddleware(jwtService)(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	jwtService := InitJWT()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/", nil)
	ctx.Request = req

	middlewares.AuthMiddleware(jwtService)(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
