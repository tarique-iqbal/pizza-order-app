package security

import (
	"errors"
	"fmt"
	"identity-service/internal/domain/auth"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtManager struct {
	secret []byte
}

func NewJWTManager(secret string) auth.JWTManager {
	return &jwtManager{secret: []byte(secret)}
}

type jwtClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *jwtManager) Generate(userID string, role string) (string, error) {
	duration := 30 * time.Minute
	if os.Getenv("APP_ENV") == "dev" {
		duration = 24 * time.Hour
	}

	expirationTime := time.Now().Add(duration)

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
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	_, err = uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user id format in token")
	}

	return &auth.UserClaims{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
