package auth

type JWTManager interface {
	Generate(userID uint, role string) (string, error)
	Parse(tokenString string) (*Claims, error)
}

type RefreshTokenManager interface {
	Generate() (string, error)
	Hash(token string) (string, error)
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}
