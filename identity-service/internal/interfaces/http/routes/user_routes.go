package routes

import (
	"identity-service/internal/interfaces/http"
	"identity-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, handler *http.UserHandler, m *middlewares.Middleware) {
	users := router.Group("/users")
	{
		authRoutes := users.Use(m.Auth)
		users.POST("", handler.Register)
		authRoutes.GET("/:id", handler.FindByID)
	}
}
