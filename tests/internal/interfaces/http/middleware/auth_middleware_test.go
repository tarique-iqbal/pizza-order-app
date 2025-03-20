package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"pizza-order-api/internal/infrastructure/security"
	"pizza-order-api/internal/interfaces/http/middlewares"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uint(1)
	token, _ := security.GenerateJWT(userID)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	ctx.Request = req

	middlewares.AuthMiddleware()(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	ctx.Request = req

	middlewares.AuthMiddleware()(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/", nil)
	ctx.Request = req

	middlewares.AuthMiddleware()(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
