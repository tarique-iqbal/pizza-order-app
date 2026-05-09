package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"identity-service/internal/container"
	"identity-service/internal/interfaces/http/routes"
	"identity-service/internal/logger"
)

func main() {
	l := logger.New()
	slog.SetDefault(l)

	c, err := container.NewAPIContainer()
	if err != nil {
		slog.Error("failed to initialize API container", "error", err)
		return
	}
	defer c.Close()

	router := gin.Default()

	handlers := &routes.Handlers{
		UserHandler: c.UserHandler,
		AuthHandler: c.AuthHandler,
	}

	routes.SetupRoutes(router, handlers, c.Middleware)

	slog.Info("starting HTTP server", "addr", ":8080")

	if err := router.Run(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
