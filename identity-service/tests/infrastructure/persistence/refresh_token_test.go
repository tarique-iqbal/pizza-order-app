package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"identity-service/internal/domain/auth"
	"identity-service/internal/infrastructure/persistence"
)

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
	claims := auth.UserClaims{
		UserID: 42,
		Role:   "owner",
	}

	err := repo.Save(ctx, hashedToken, claims, 60)
	require.NoError(t, err)

	foundClaims, err := repo.Find(ctx, hashedToken)
	require.NoError(t, err)

	assert.Equal(t, claims, foundClaims)
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
	claims := auth.UserClaims{
		UserID: 99,
		Role:   "owner",
	}

	err := repo.Save(ctx, hashedToken, claims, 60)
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
	claims := auth.UserClaims{
		UserID: 7,
		Role:   "owner",
	}

	err := repo.Save(ctx, hashedToken, claims, 1)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	_, err = repo.Find(ctx, hashedToken)

	require.Error(t, err)
	assert.Equal(t, "refresh token not found", err.Error())
}
