package container

import (
	aRestaurant "api-service/internal/application/restaurant"
	"api-service/internal/infrastructure/db"
	"api-service/internal/infrastructure/geocoder"
	"api-service/internal/infrastructure/messaging"
	"api-service/internal/infrastructure/persistence"
	"api-service/internal/infrastructure/security"
	"api-service/internal/interfaces/http"
	"api-service/internal/interfaces/http/middlewares"
	"os"

	"gorm.io/gorm"
)

type Container struct {
	RestaurantHandler *http.RestaurantHandler
	PizzaSizeHandler  *http.PizzaSizeHandler
	DB                *gorm.DB
	Publisher         *messaging.RabbitMQPublisher
	Middleware        *middlewares.Middleware
}

func NewContainer() (*Container, error) {
	database, _ := db.InitDB()

	amqpURL := os.Getenv("RABBITMQ_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	opencageApiKey := os.Getenv("OPENCAGE_API_KEY")

	publisher := messaging.NewRabbitMQPublisher(amqpURL)
	jwtService := security.NewJWTService(jwtSecret)
	middleware := middlewares.NewMiddleware(jwtService)

	restaurantRepo := persistence.NewRestaurantRepository(database)

	// restaurant
	geocoderService := geocoder.NewOpenCageService(opencageApiKey)
	restaurantAddressRepo := persistence.NewRestaurantAddressRepository(database)
	createRestaurantUC := aRestaurant.NewCreateRestaurantUseCase(database, geocoderService, restaurantRepo, restaurantAddressRepo)
	restaurantHandler := http.NewRestaurantHandler(createRestaurantUC)

	// pizza-sizes
	pizzaSizeRepo := persistence.NewPizzaSizeRepository(database)
	createPizzaSizeUC := aRestaurant.NewCreatePizzaSizeUseCase(pizzaSizeRepo, restaurantRepo)
	pizzaSizeHandler := http.NewPizzaSizeHandler(createPizzaSizeUC)

	return &Container{
		RestaurantHandler: restaurantHandler,
		PizzaSizeHandler:  pizzaSizeHandler,
		DB:                database,
		Publisher:         publisher,
		Middleware:        middleware,
	}, nil
}

func (c *Container) Close() {
	db, _ := c.DB.DB()
	db.Close()

	c.Publisher.Close()
}
