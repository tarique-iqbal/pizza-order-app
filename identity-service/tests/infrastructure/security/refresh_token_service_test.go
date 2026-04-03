package security_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"identity-service/internal/infrastructure/security"
)

func TestGenerate_ReturnsValidToken(t *testing.T) {
	service := security.NewRefreshTokenService()

	token, err := service.Generate()

	require.NoError(t, err)
	require.NotEmpty(t, token)

	assert.Len(t, token, 64)

	_, err = hex.DecodeString(token)
	assert.NoError(t, err)
}

func TestGenerate_UniqueTokens(t *testing.T) {
	service := security.NewRefreshTokenService()

	token1, err1 := service.Generate()
	token2, err2 := service.Generate()

	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.NotEqual(t, token1, token2)
}

func TestHash_Deterministic(t *testing.T) {
	service := security.NewRefreshTokenService()

	input := "my-refresh-token"

	hash1, err1 := service.Hash(input)
	hash2, err2 := service.Hash(input)

	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.Equal(t, hash1, hash2)
}

func TestHash_DifferentInputs(t *testing.T) {
	service := security.NewRefreshTokenService()

	hash1, err1 := service.Hash("token-1")
	hash2, err2 := service.Hash("token-2")

	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.NotEqual(t, hash1, hash2)
}

func TestHash_OutputLength(t *testing.T) {
	service := security.NewRefreshTokenService()

	hash, err := service.Hash("some-token")

	require.NoError(t, err)

	assert.Len(t, hash, 64)
}
