package main

import (
	"log"
	"restaurant-service/internal/container"
	"restaurant-service/internal/interfaces/http/routes"

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
		RestaurantHandler: iocContainer.RestaurantHandler,
		PizzaSizeHandler:  iocContainer.PizzaSizeHandler,
	}

	routes.SetupRoutes(router, handlers, iocContainer.Middleware)

	router.Run(":8080")
}
