package main

import (
	"log"
	"pizza-order-api/container"
	"pizza-order-api/internal/interfaces/http/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	iocContainer, err := container.NewContainer()
	if err != nil {
		log.Fatal("Failed to initialize container:", err)
	}
	defer iocContainer.Close()

	router := gin.Default()
	handlers := &routes.Handlers{
		UserHandler:       iocContainer.UserHandler,
		AuthHandler:       iocContainer.AuthHandler,
		RestaurantHandler: iocContainer.RestaurantHandler,
	}

	routes.SetupRoutes(router, handlers)

	router.Run(":8080")
}
