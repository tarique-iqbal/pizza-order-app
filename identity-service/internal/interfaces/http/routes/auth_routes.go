package routes

import (
	"identity-service/internal/interfaces/http"
	"identity-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, handler *http.AuthHandler, m *middlewares.Middleware) {
	auth := router.Group("/auth")

	auth.POST("/email/verify", handler.CreateEmailVerification)
	auth.POST("/login", handler.Login)

	protected := auth.Group("")
	protected.Use(m.Auth)

	protected.POST("/refresh", handler.Refresh)
	protected.POST("/logout", handler.Logout)
}
