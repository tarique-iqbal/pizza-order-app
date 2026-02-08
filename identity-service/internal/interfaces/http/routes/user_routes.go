package routes

import (
	"identity-service/internal/interfaces/http"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userHandler *http.UserHandler) {
	users := router.Group("/users")
	{
		users.POST("", userHandler.Create)
	}
}
