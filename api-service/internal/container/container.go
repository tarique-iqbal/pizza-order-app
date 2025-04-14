package container

import (
	"api-service/internal/application/auth"
	"api-service/internal/application/restaurant"
	"api-service/internal/application/user"
	"api-service/internal/infrastructure/db"
	"api-service/internal/infrastructure/messaging"
	"api-service/internal/infrastructure/persistence"
	"api-service/internal/infrastructure/validator"
	"api-service/internal/interfaces/http"
	"os"

	"gorm.io/gorm"
)

type Container struct {
	UserHandler       *http.UserHandler
	AuthHandler       *http.AuthHandler
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

	// user
	createUserUseCase := user.NewCreateUserUseCase(userRepo, publisher)
	userUseCases := &http.UserUseCases{
		CreateUser:      createUserUseCase,
		CustomValidator: customValidator,
	}
	userHandler := http.NewUserHandler(userUseCases)

	// auth
	signInUseCase := auth.NewSignInUseCase(userRepo)
	authUseCases := &http.AuthUseCases{
		SignIn: signInUseCase,
	}
	authHandler := http.NewAuthHandler(authUseCases)

	// restaurant
	createRestaurantUseCase := restaurant.NewCreateRestaurantUseCase(restaurantRepo)
	restaurantUseCase := &http.RestaurantUseCases{
		CreateRestaurant: createRestaurantUseCase,
		CustomValidator:  customValidator,
	}
	restaurantHandler := http.NewRestaurantHandler(restaurantUseCase)

	return &Container{
		UserHandler:       userHandler,
		AuthHandler:       authHandler,
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
