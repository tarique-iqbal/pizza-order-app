package auth

type JWTService interface {
	Generate(userID uint, role string) (string, error)
	Parse(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}
