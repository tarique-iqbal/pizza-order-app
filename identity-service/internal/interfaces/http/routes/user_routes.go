package routes

import (
	"identity-service/internal/interfaces/http"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, handler *http.UserHandler) {
	users := router.Group("/users")
	{
		users.POST("", handler.Register)
	}
}
