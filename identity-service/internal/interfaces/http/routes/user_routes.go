package routes

import (
	"identity-service/internal/interfaces/http"
	"identity-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, handler *http.UserHandler, m *middlewares.Middleware) {
	users := router.Group("/users")

	users.POST("/owners", handler.RegisterOwner)
	users.POST("/customers", handler.RegisterCustomer)

	protected := users.Group("")
	protected.Use(m.Auth)

	protected.GET("/:id", handler.FindByID)
}
