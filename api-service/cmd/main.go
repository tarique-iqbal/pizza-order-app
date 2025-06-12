package main

import (
	"api-service/internal/container"
	"api-service/internal/interfaces/http/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env != "docker" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
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
		PizzaSizeHandler:  iocContainer.PizzaSizeHandler,
	}

	routes.SetupRoutes(router, handlers, iocContainer.Middleware)

	router.Run(":8080")
}
