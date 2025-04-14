package routes

import (
	"api-service/internal/interfaces/http"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userHandler *http.UserHandler) {
	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("", userHandler.Create)
		}
	}
}
