package container

import (
	"os"
	aRestaurant "restaurant-service/internal/application/restaurant"
	"restaurant-service/internal/infrastructure/db"
	"restaurant-service/internal/infrastructure/geocoder"
	"restaurant-service/internal/infrastructure/messaging"
	"restaurant-service/internal/infrastructure/persistence"
	"restaurant-service/internal/infrastructure/security"
	"restaurant-service/internal/interfaces/http"
	"restaurant-service/internal/interfaces/http/middlewares"

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
