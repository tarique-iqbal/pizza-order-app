package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type RefreshTokenService struct {
}

func NewRefreshTokenService() *RefreshTokenService {
	return &RefreshTokenService{}
}

func (s *RefreshTokenService) Generate() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *RefreshTokenService) Hash(token string) (string, error) {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:]), nil
}
