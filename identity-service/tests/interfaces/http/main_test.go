package http_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"identity-service/internal/infrastructure/redis"
	"identity-service/tests/infrastructure/db"
)

type TestStorage struct {
	DB    *gorm.DB
	Redis *goredis.Client
}

var ts *TestStorage

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	cfg := redis.Config{
		Addr: os.Getenv("REDIS_ADDR"),
		DB:   1,
	}

	tdb := db.SetupTestDB()
	trc, _ := redis.InitRedis(cfg)

	ts = &TestStorage{
		DB:    tdb,
		Redis: trc,
	}

	code := m.Run()
	os.Exit(code)
}

func testStorage() *TestStorage {
	return ts
}

func truncateTables(tdb *gorm.DB) {
	tdb.Exec("TRUNCATE TABLE users, email_verifications RESTART IDENTITY CASCADE")
}

func MockAuthMiddleware(userID int, role string) gin.HandlerFunc {
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
