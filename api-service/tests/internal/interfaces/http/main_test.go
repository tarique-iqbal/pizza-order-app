package http_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-service/tests/internal/infrastructure/db"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testDB = db.SetupTestDB()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func MockAuthMiddleware(expectedRole string) gin.HandlerFunc {
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

		ctx.Set("userID", uint(1))
		ctx.Set("role", expectedRole)
		ctx.Next()
	}
}
