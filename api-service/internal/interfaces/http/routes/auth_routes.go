package routes

import (
	"api-service/internal/interfaces/http"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authHandler *http.AuthHandler) {
	api := router.Group("/api")
	{
		users := api.Group("/auth")
		{
			users.POST("/email-verification", authHandler.CreateEmailVerification)
			users.POST("/signin", authHandler.SignIn)
		}
	}
}
