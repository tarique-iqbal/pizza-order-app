package persistence_test

import (
	"context"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
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

func flushRedis(t *testing.T, client *goredis.Client) {
	err := client.FlushDB(context.Background()).Err()
	require.NoError(t, err)
}

func generateUUID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id
}
