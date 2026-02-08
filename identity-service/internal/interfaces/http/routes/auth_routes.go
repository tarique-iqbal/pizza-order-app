package routes

import (
	"identity-service/internal/interfaces/http"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authHandler *http.AuthHandler) {
	users := router.Group("/auth")
	{
		users.POST("/email-verification", authHandler.CreateEmailVerification)
		users.POST("/signin", authHandler.SignIn)
	}
}
