package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"identity-service/internal/domain/auth"
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
	claims auth.UserClaims,
	ttlSeconds int64,
) error {
	key := "refresh:" + hashedToken
	bytes, err := json.Marshal(claims)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, bytes, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *refreshTokenRepo) Find(
	ctx context.Context,
	hashedToken string,
) (auth.UserClaims, error) {
	key := "refresh:" + hashedToken

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return auth.UserClaims{}, errors.New("refresh token not found")
	}

	var claims auth.UserClaims

	if err := json.Unmarshal([]byte(val), &claims); err != nil {
		return auth.UserClaims{}, errors.New("invalid token data")
	}

	return claims, nil
}

func (r *refreshTokenRepo) Delete(
	ctx context.Context,
	hashedToken string,
) error {
	key := "refresh:" + hashedToken
	return r.client.Del(ctx, key).Err()
}
