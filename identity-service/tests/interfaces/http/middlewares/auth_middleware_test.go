package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"identity-service/internal/domain/auth"
	"identity-service/internal/infrastructure/security"
	"identity-service/internal/interfaces/http/middlewares"
	"identity-service/tests/testutil"
)

func InitJWT() auth.JWTManager {
	gin.SetMode(gin.TestMode)
	return security.NewJWTManager("TestSecretKey")
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtManager := InitJWT()

	userID := testutil.MustNewID()
	role := "customer"

	token, _ := jwtManager.Generate(userID.String(), role)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	ctx.Request = req

	middlewares.AuthMiddleware(jwtManager)(ctx)
	ctxUserID := ctx.MustGet(middlewares.CtxUserID)
	ctxUserRole := ctx.MustGet(middlewares.CtxUserRole)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, userID.String(), ctxUserID)
	assert.Equal(t, role, ctxUserRole)
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
