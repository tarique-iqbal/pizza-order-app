package middlewares

import (
	"api-service/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Auth gin.HandlerFunc
}

func NewMiddleware(jwt auth.JWTService) *Middleware {
	return &Middleware{
		Auth: AuthMiddleware(jwt),
	}
}
