package routes

import (
	"identity-service/internal/interfaces/http"
	"identity-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, handler *http.AuthHandler, m *middlewares.Middleware) {
	routes := router.Group("/auth")

	routes.POST("/email-verification", handler.CreateEmailVerification)
	routes.POST("/login", handler.Login)

	authRoutes := routes.Group("/")
	authRoutes.Use(m.Auth)
	{
		authRoutes.POST("/refresh", handler.Refresh)
		authRoutes.POST("/logout", handler.Logout)
	}
}
