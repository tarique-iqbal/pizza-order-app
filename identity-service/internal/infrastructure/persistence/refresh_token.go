package persistence

import (
	"context"
	"errors"
	"identity-service/internal/domain/auth"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type refreshTokenRepo struct {
	client *redis.Client
}

func NewRefreshTokenRepository(client *redis.Client) auth.RefreshTokenRepository {
	return &refreshTokenRepo{client: client}
}

func (r *refreshTokenRepo) Save(
	ctx context.Context,
	hashedToken string,
	userID int,
	ttlSeconds int64,
) error {
	key := "refresh:" + hashedToken
	value := strconv.Itoa(userID)

	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *refreshTokenRepo) Find(
	ctx context.Context,
	hashedToken string,
) (int, error) {
	key := "refresh:" + hashedToken

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return 0, errors.New("refresh token not found")
	}

	userID, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.New("invalid stored value")
	}

	return userID, nil
}

func (r *refreshTokenRepo) Delete(
	ctx context.Context,
	hashedToken string,
) error {
	key := "refresh:" + hashedToken
	return r.client.Del(ctx, key).Err()
}
