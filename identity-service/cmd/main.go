package main

import (
	"identity-service/internal/container"
	"identity-service/internal/interfaces/http/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	iocContainer, err := container.NewContainer()
	if err != nil {
		log.Fatal("Failed to initialize container:", err)
	}
	defer iocContainer.Close()

	router := gin.Default()
	handlers := &routes.Handlers{
		UserHandler: iocContainer.UserHandler,
		AuthHandler: iocContainer.AuthHandler,
	}

	routes.SetupRoutes(router, handlers, iocContainer.Middleware)

	router.Run(":8080")
}
