package container

import (
	"log"
	"pizza-order-api/internal/application/restaurant"
	"pizza-order-api/internal/application/user"
	"pizza-order-api/internal/infrastructure/db"
	"pizza-order-api/internal/infrastructure/persistence"
	"pizza-order-api/internal/infrastructure/validator"
	"pizza-order-api/internal/interfaces/http"

	"gorm.io/gorm"
)

type Container struct {
	UserHandler       *http.UserHandler
	RestaurantHandler *http.RestaurantHandler
	DB                *gorm.DB
}

func NewContainer() (*Container, error) {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("Could not connect to database:", err)
		return nil, err
	}

	userRepo := persistence.NewUserRepository(database)
	restaurantRepo := persistence.NewRestaurantRepository(database)

	customValidator := validator.NewCustomValidator(userRepo, restaurantRepo)

	createUserUseCase := user.NewCreateUserUseCase(userRepo)
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
	}, nil
}

func (c *Container) Close() {
	db, _ := c.DB.DB()
	db.Close()
}
