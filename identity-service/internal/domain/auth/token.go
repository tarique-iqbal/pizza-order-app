package auth

type JWTManager interface {
	Generate(userID int, role string) (string, error)
	Parse(tokenString string) (*Claims, error)
}

type RefreshTokenManager interface {
	Generate() (string, error)
	Hash(token string) (string, error)
}

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
}
