package security

import (
	"errors"
	"identity-service/internal/domain/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtManager struct {
	secret []byte
}

func NewJWTManager(secret string) auth.JWTManager {
	return &jwtManager{secret: []byte(secret)}
}

type jwtClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *jwtManager) Generate(userID int, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwtClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *jwtManager) Parse(tokenString string) (*auth.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &auth.UserClaims{UserID: claims.UserID, Role: claims.Role}, nil
	}

	return nil, errors.New("invalid token")
}
