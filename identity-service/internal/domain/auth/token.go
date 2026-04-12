package auth

type JWTManager interface {
	Generate(userID int, role string) (string, error)
	Parse(tokenString string) (*UserClaims, error)
}

type RefreshTokenManager interface {
	Generate() (string, error)
	Hash(token string) string
}

type UserClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
}
