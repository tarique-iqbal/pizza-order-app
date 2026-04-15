package routes

import (
	"identity-service/internal/interfaces/http"
	"identity-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, handler *http.AuthHandler, m *middlewares.Middleware) {
	users := router.Group("/auth")
	{
		authRoutes := users.Use(m.Auth)
		users.POST("/email-verification", handler.CreateEmailVerification)
		users.POST("/login", handler.Login)
		users.POST("/refresh", handler.Refresh)
		authRoutes.POST("/logout", handler.Logout)
	}
}
