package http_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"api-service/internal/domain/auth"
	"api-service/internal/domain/restaurant"
	"api-service/internal/domain/user"
	"api-service/tests/infrastructure/db"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testDB = db.SetupTestDB()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func MockAuthMiddleware(userID uint, role string) gin.HandlerFunc {
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

		ctx.Set("userID", userID)
		ctx.Set("role", role)
		ctx.Next()
	}
}

func resetTables(t *testing.T) {
	tables := []string{
		auth.EmailVerification{}.TableName(),
		user.User{}.TableName(),
		restaurant.Restaurant{}.TableName(),
		restaurant.PizzaSize{}.TableName(),
	}

	for _, table := range tables {
		err := testDB.Exec("DELETE FROM " + table).Error
		require.NoError(t, err)
	}
}
