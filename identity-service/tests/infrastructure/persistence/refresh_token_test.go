package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"identity-service/internal/infrastructure/persistence"
)

func flushRedis(t *testing.T, client *redis.Client) {
	err := client.FlushDB(context.Background()).Err()
	require.NoError(t, err)
}

func TestRedisConnection(t *testing.T) {
	ts := testStorage()
	err := ts.Redis.Ping(context.Background()).Err()
	require.NoError(t, err)
}

func TestRefreshTokenRepository_SaveAndFind(t *testing.T) {
	ctx := context.Background()

	ts := testStorage()
	flushRedis(t, ts.Redis)

	repo := persistence.NewRefreshTokenRepository(ts.Redis)

	hashedToken := "test-token"
	userID := 42

	err := repo.Save(ctx, hashedToken, userID, 60)
	require.NoError(t, err)

	foundUserID, err := repo.Find(ctx, hashedToken)
	require.NoError(t, err)

	assert.Equal(t, userID, foundUserID)
}

func TestRefreshTokenRepository_Find_NotFound(t *testing.T) {
	ctx := context.Background()

	ts := testStorage()
	flushRedis(t, ts.Redis)

	repo := persistence.NewRefreshTokenRepository(ts.Redis)

	_, err := repo.Find(ctx, "non-existing-token")

	require.Error(t, err)
	assert.Equal(t, "refresh token not found", err.Error())
}

func TestRefreshTokenRepository_Delete(t *testing.T) {
	ctx := context.Background()

	ts := testStorage()
	flushRedis(t, ts.Redis)

	repo := persistence.NewRefreshTokenRepository(ts.Redis)

	hashedToken := "delete-token"
	userID := 99

	err := repo.Save(ctx, hashedToken, userID, 60)
	require.NoError(t, err)

	err = repo.Delete(ctx, hashedToken)
	require.NoError(t, err)

	_, err = repo.Find(ctx, hashedToken)
	require.Error(t, err)
}

func TestRefreshTokenRepository_TTLExpiry(t *testing.T) {
	ctx := context.Background()

	ts := testStorage()
	flushRedis(t, ts.Redis)

	repo := persistence.NewRefreshTokenRepository(ts.Redis)

	hashedToken := "ttl-token"
	userID := 7

	err := repo.Save(ctx, hashedToken, userID, 1)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	_, err = repo.Find(ctx, hashedToken)

	require.Error(t, err)
	assert.Equal(t, "refresh token not found", err.Error())
}
