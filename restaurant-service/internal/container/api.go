package container

import (
	"os"

	resapp "restaurant-service/internal/application/restaurant"
	"restaurant-service/internal/infrastructure/geocoder"
	"restaurant-service/internal/infrastructure/messaging"
	"restaurant-service/internal/infrastructure/persistence"
	"restaurant-service/internal/interfaces/http"
	"restaurant-service/internal/interfaces/http/middlewares"
)

type APIContainer struct {
	*Shared
	Middleware        *middlewares.Middleware
	Publisher         *messaging.RabbitMQPublisher
	RestaurantHandler *http.RestaurantHandler
	PizzaSizeHandler  *http.PizzaSizeHandler
}

func NewAPIContainer() (*APIContainer, error) {
	base, err := NewShared()
	if err != nil {
		return nil, err
	}

	opencageApiKey := os.Getenv("OPENCAGE_API_KEY")

	publisher := messaging.NewRabbitMQPublisher(base.AMQPURL)
	middleware := middlewares.NewMiddleware()

	restaurantRepo := persistence.NewRestaurantRepository(base.DB)

	// restaurant
	geocoder := geocoder.NewOpenCageGeocoder(opencageApiKey)
	restaurantAddressRepo := persistence.NewRestaurantAddressRepository(base.DB)
	createRestaurant := resapp.NewCreateRestaurant(base.DB, geocoder, restaurantRepo, restaurantAddressRepo)
	restaurantHandler := http.NewRestaurantHandler(createRestaurant)

	// pizza-sizes
	pizzaSizeRepo := persistence.NewPizzaSizeRepository(base.DB)
	createPizzaSize := resapp.NewCreatePizzaSize(pizzaSizeRepo, restaurantRepo)
	pizzaSizeHandler := http.NewPizzaSizeHandler(createPizzaSize)

	return &APIContainer{
		Shared:            base,
		Middleware:        middleware,
		Publisher:         publisher,
		RestaurantHandler: restaurantHandler,
		PizzaSizeHandler:  pizzaSizeHandler,
	}, nil
}

func (c *APIContainer) Close() {
	if c.DB != nil {
		db, err := c.DB.DB()
		if err == nil {
			_ = db.Close()
		}
	}

	if c.Publisher != nil {
		c.Publisher.Close()
	}
}
