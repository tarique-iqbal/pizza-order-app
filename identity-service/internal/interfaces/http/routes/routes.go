package routes

import (
	"identity-service/internal/interfaces/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	UserHandler *http.UserHandler
	AuthHandler *http.AuthHandler
}

func SetupRoutes(router *gin.Engine, h *Handlers) {
	SetupUserRoutes(router, h.UserHandler)
	SetupAuthRoutes(router, h.AuthHandler)
}
