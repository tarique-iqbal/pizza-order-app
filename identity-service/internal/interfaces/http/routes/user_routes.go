package routes

import (
	"identity-service/internal/interfaces/http"
	"identity-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, handler *http.UserHandler, m *middlewares.Middleware) {
	routes := router.Group("/users")

	routes.POST("", handler.Register)

	authRoutes := routes.Group("/")
	authRoutes.Use(m.Auth)
	{
		authRoutes.GET("/:id", handler.FindByID)
	}
}
