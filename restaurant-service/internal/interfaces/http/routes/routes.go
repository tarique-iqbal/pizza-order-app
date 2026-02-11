package routes

import (
	"restaurant-service/internal/interfaces/http"
	"restaurant-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	RestaurantHandler *http.RestaurantHandler
	PizzaSizeHandler  *http.PizzaSizeHandler
}

func SetupRoutes(router *gin.Engine, h *Handlers, m *middlewares.Middleware) {
	SetupRestaurantRoutes(router, h.RestaurantHandler, m)
	SetupPizzaSizeRoutes(router, h.PizzaSizeHandler, m)
}
