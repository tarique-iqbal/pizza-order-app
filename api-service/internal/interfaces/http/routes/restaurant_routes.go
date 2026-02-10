package routes

import (
	"api-service/internal/interfaces/http"
	"api-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRestaurantRoutes(router *gin.Engine, rh *http.RestaurantHandler, m *middlewares.Middleware) {
	restaurants := router.Group("/restaurants")
	{
		authRoutes := restaurants.Use(m.Auth, m.EnsureOwner)
		authRoutes.POST("", rh.Create)
	}
}
