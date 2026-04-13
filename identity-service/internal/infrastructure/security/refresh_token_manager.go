package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"identity-service/internal/domain/auth"
)

type RefreshTokenManager struct {
}

func NewRefreshTokenManager() auth.RefreshTokenManager {
	return &RefreshTokenManager{}
}

func (s *RefreshTokenManager) Generate() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *RefreshTokenManager) Hash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
