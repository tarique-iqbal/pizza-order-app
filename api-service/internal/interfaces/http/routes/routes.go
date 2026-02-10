package routes

import (
	"api-service/internal/interfaces/http"
	"api-service/internal/interfaces/http/middlewares"

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
