package routes

import (
	"pizza-order-api/internal/interfaces/http"
	"pizza-order-api/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRestaurantRoutes(router *gin.Engine, rh *http.RestaurantHandler) {
	api := router.Group("/api")
	{
		restaurants := api.Group("/restaurants")
		{
			authRoutes := restaurants.Use(middlewares.AuthMiddleware())
			authRoutes.POST("", rh.Create)
		}
	}
}
