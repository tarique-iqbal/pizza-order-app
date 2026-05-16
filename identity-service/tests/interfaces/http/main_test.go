package http_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"identity-service/internal/interfaces/http/middlewares"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	code := m.Run()
	os.Exit(code)
}

func MockAuthMiddleware(userID string, role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		if authHeader != "Bearer mock-valid-token" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set(middlewares.CtxUserID, userID)
		ctx.Set(middlewares.CtxUserRole, role)

		ctx.Next()
	}
}
