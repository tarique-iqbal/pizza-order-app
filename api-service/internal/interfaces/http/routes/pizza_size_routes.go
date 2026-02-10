package routes

import (
	"api-service/internal/interfaces/http"
	"api-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupPizzaSizeRoutes(router *gin.Engine, psh *http.PizzaSizeHandler, m *middlewares.Middleware) {
	restaurants := router.Group("/restaurants")
	{
		authRoutes := restaurants.Use(m.Auth, m.EnsureOwner)
		authRoutes.POST("/:id/pizza-sizes", psh.Create)
	}
}
