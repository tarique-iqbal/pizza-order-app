package container

import (
	"os"
	"pizza-order-api/internal/application/restaurant"
	"pizza-order-api/internal/application/user"
	"pizza-order-api/internal/infrastructure/db"
	"pizza-order-api/internal/infrastructure/messaging"
	"pizza-order-api/internal/infrastructure/persistence"
	"pizza-order-api/internal/infrastructure/validator"
	"pizza-order-api/internal/interfaces/http"

	"gorm.io/gorm"
)

type Container struct {
	UserHandler       *http.UserHandler
	RestaurantHandler *http.RestaurantHandler
	DB                *gorm.DB
	Publisher         *messaging.RabbitMQPublisher
}

func NewContainer() (*Container, error) {
	database, _ := db.InitDB()

	amqpURL := os.Getenv("RABBITMQ_URL")
	publisher := messaging.NewRabbitMQPublisher(amqpURL)

	userRepo := persistence.NewUserRepository(database)
	restaurantRepo := persistence.NewRestaurantRepository(database)

	customValidator := validator.NewCustomValidator(userRepo, restaurantRepo)

	createUserUseCase := user.NewCreateUserUseCase(userRepo, publisher)
	signInUserUseCase := user.NewSignInUserUseCase(userRepo)

	userUseCases := &http.UserUseCases{
		CreateUser:      createUserUseCase,
		SignIn:          signInUserUseCase,
		CustomValidator: customValidator,
	}
	userHandler := http.NewUserHandler(userUseCases)

	createRestaurantUseCase := restaurant.NewCreateRestaurantUseCase(restaurantRepo)

	restaurantUseCase := &http.RestaurantUseCases{
		CreateRestaurant: createRestaurantUseCase,
		CustomValidator:  customValidator,
	}
	restaurantHandler := http.NewRestaurantHandler(restaurantUseCase)

	return &Container{
		UserHandler:       userHandler,
		RestaurantHandler: restaurantHandler,
		DB:                database,
		Publisher:         publisher,
	}, nil
}

func (c *Container) Close() {
	db, _ := c.DB.DB()
	db.Close()

	c.Publisher.Close()
}
