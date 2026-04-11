package middlewares

import (
	"identity-service/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Auth gin.HandlerFunc
}

func NewMiddleware(jwt auth.JWTManager) *Middleware {
	return &Middleware{
		Auth: AuthMiddleware(jwt),
	}
}
