package routes

import (
	"identity-service/internal/interfaces/http"
	"identity-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, handler *http.UserHandler, m *middlewares.Middleware) {
	router.POST("/owners", handler.RegisterOwner)
	router.POST("/customers", handler.RegisterCustomer)

	authRoutes := router.Group("/users")
	authRoutes.Use(m.Auth)
	{
		authRoutes.GET("/:id", handler.FindByID)
	}
}
