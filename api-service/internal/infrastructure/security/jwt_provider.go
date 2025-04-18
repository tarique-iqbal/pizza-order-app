package security

import (
	"api-service/internal/domain/auth"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secret []byte
}

func NewJWTService(secret string) auth.JWTService {
	return &jwtService{secret: []byte(secret)}
}

type jwtClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func (j *jwtService) GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *jwtService) ParseToken(tokenString string) (*auth.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &auth.Claims{UserID: claims.UserID}, nil
	}

	return nil, errors.New("invalid token")
}
