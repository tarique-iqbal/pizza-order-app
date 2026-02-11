package middlewares

import (
	"restaurant-service/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Auth        gin.HandlerFunc
	EnsureOwner gin.HandlerFunc
}

func NewMiddleware(jwt auth.JWTService) *Middleware {
	return &Middleware{
		Auth:        AuthMiddleware(jwt),
		EnsureOwner: RequireRole("owner"),
	}
}
