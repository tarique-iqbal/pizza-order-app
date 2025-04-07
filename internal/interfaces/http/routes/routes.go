package routes

import (
	"pizza-order-api/internal/interfaces/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	UserHandler       *http.UserHandler
	AuthHandler       *http.AuthHandler
	RestaurantHandler *http.RestaurantHandler
}

func SetupRoutes(router *gin.Engine, h *Handlers) {
	SetupUserRoutes(router, h.UserHandler)
	SetupAuthRoutes(router, h.AuthHandler)
	SetupRestaurantRoutes(router, h.RestaurantHandler)
}
