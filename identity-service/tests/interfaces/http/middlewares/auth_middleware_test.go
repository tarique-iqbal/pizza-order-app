package middlewares_test

import (
	"identity-service/internal/domain/auth"
	"identity-service/internal/infrastructure/security"
	"identity-service/internal/interfaces/http/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func InitJWT() auth.JWTManager {
	gin.SetMode(gin.TestMode)
	return security.NewJWTManager("TestSecretKey")
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtManager := InitJWT()

	userID, _ := uuid.NewV7()
	role := "customer"

	token, _ := jwtManager.Generate(userID.String(), role)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	ctx.Request = req

	middlewares.AuthMiddleware(jwtManager)(ctx)
	ctxUserID := ctx.MustGet("userID")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, userID.String(), ctxUserID)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwtManager := InitJWT()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	ctx.Request = req

	middlewares.AuthMiddleware(jwtManager)(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	jwtManager := InitJWT()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/", nil)
	ctx.Request = req

	middlewares.AuthMiddleware(jwtManager)(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
