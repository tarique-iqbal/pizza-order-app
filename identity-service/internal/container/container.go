package container

import (
	aAuth "identity-service/internal/application/auth"
	"identity-service/internal/application/user"
	iAuth "identity-service/internal/infrastructure/auth"
	"identity-service/internal/infrastructure/db"
	"identity-service/internal/infrastructure/messaging"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	"identity-service/internal/interfaces/http"
	"identity-service/internal/interfaces/http/middlewares"
	"os"

	"gorm.io/gorm"
)

type Container struct {
	UserHandler *http.UserHandler
	AuthHandler *http.AuthHandler
	DB          *gorm.DB
	Publisher   *messaging.RabbitMQPublisher
	Middleware  *middlewares.Middleware
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

	codeVerifier := iAuth.NewCodeVerificationService(emailVerificationRepo)

	// user
	createUserUC := user.NewCreateUserUseCase(codeVerifier, userRepo, hasher, publisher)
	userHandler := http.NewUserHandler(createUserUC)

	// auth
	signInUC := aAuth.NewSignInUseCase(userRepo, hasher, jwtService)
	createEmailVerificationUC := aAuth.NewCreateEmailVerificationUseCase(emailVerificationRepo, otp, publisher)
	authHandler := http.NewAuthHandler(signInUC, createEmailVerificationUC)

	return &Container{
		UserHandler: userHandler,
		AuthHandler: authHandler,
		DB:          database,
		Publisher:   publisher,
		Middleware:  middleware,
	}, nil
}

func (c *Container) Close() {
	db, _ := c.DB.DB()
	db.Close()

	c.Publisher.Close()
}
