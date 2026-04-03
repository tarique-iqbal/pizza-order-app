package auth

type RefreshTokenService interface {
	Generate() (string, error)
	Hash(token string) (string, error)
}
