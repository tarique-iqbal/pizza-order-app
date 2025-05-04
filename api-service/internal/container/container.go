package container

import (
	aAuth "api-service/internal/application/auth"
	aRestaurant "api-service/internal/application/restaurant"
	"api-service/internal/application/user"
	iAuth "api-service/internal/infrastructure/auth"
	"api-service/internal/infrastructure/db"
	"api-service/internal/infrastructure/messaging"
	"api-service/internal/infrastructure/persistence"
	"api-service/internal/infrastructure/security"
	"api-service/internal/infrastructure/validator"
	"api-service/internal/interfaces/http"
	"api-service/internal/interfaces/http/middlewares"
	"os"

	"gorm.io/gorm"
)

type Container struct {
	UserHandler       *http.UserHandler
	AuthHandler       *http.AuthHandler
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

	publisher := messaging.NewRabbitMQPublisher(amqpURL)
	hasher := security.NewPasswordHasher()
	jwtService := security.NewJWTService(jwtSecret)
	middleware := middlewares.NewMiddleware(jwtService)
	otp := security.NewSixDigitOTPGenerator()

	emailVerificationRepo := persistence.NewEmailVerificationRepository(database)
	userRepo := persistence.NewUserRepository(database)
	restaurantRepo := persistence.NewRestaurantRepository(database)
	customValidator := validator.NewCustomValidator(userRepo, restaurantRepo)

	codeVerifier := iAuth.NewCodeVerificationService(emailVerificationRepo)

	// user
	createUserUseCase := user.NewCreateUserUseCase(codeVerifier, userRepo, hasher, publisher)
	userUseCases := &http.UserUseCases{
		CreateUser:      createUserUseCase,
		CustomValidator: customValidator,
	}
	userHandler := http.NewUserHandler(userUseCases)

	// auth
	signInUseCase := aAuth.NewSignInUseCase(userRepo, hasher, jwtService)
	createEmailVerificationUseCase := aAuth.NewCreateEmailVerificationUseCase(emailVerificationRepo, otp, publisher)
	authUseCases := &http.AuthUseCases{
		SignIn:                  signInUseCase,
		CreateEmailVerification: createEmailVerificationUseCase,
	}
	authHandler := http.NewAuthHandler(authUseCases)

	// restaurant
	createRestaurantUseCase := aRestaurant.NewCreateRestaurantUseCase(restaurantRepo)
	restaurantUseCase := &http.RestaurantUseCases{
		CreateRestaurant: createRestaurantUseCase,
		CustomValidator:  customValidator,
	}
	restaurantHandler := http.NewRestaurantHandler(restaurantUseCase)

	// pizza-sizes
	pizzaSizeRepo := persistence.NewPizzaSizeRepository(database)
	createPizzaSizeUC := aRestaurant.NewCreatePizzaSizeUseCase(pizzaSizeRepo, restaurantRepo)
	pizzaSizeHandler := http.NewPizzaSizeHandler(createPizzaSizeUC)

	return &Container{
		UserHandler:       userHandler,
		AuthHandler:       authHandler,
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
