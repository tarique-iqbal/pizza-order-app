package http_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"pizza-order-api/tests/internal/infrastructure/db"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	testDB = db.SetupTestDB()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func MockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if authHeader != "Bearer mock-valid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", uint(1))
		c.Next()
	}
}
