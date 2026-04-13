package routes

import (
	"identity-service/internal/interfaces/http"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, handler *http.AuthHandler) {
	users := router.Group("/auth")
	{
		users.POST("/email-verification", handler.CreateEmailVerification)
		users.POST("/login", handler.Login)
		users.POST("/refresh", handler.Refresh)
	}
}
