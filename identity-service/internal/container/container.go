package container

import (
	aAuth "identity-service/internal/application/auth"
	"identity-service/internal/application/user"
	iAuth "identity-service/internal/infrastructure/auth"
	"identity-service/internal/infrastructure/db"
	"identity-service/internal/infrastructure/messaging"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/redis"
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
	amqpURL := os.Getenv("RABBITMQ_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	cfg := redis.Config{
		Addr: os.Getenv("REDIS_ADDR"),
	}

	database, _ := db.InitDB()
	rc, _ := redis.InitRedis(cfg)

	publisher := messaging.NewRabbitMQPublisher(amqpURL)
	hasher := security.NewPasswordHasher()
	jwtManager := security.NewJWTManager(jwtSecret)
	refreshTokenManager := security.NewRefreshTokenManager()
	middleware := middlewares.NewMiddleware(jwtManager)
	otp := security.NewOTPGenerator()

	refreshTokenRepo := persistence.NewRefreshTokenRepository(rc)
	emailVerificationRepo := persistence.NewEmailVerificationRepository(database)
	userRepo := persistence.NewUserRepository(database)

	codeVerifier := iAuth.NewEmailVerifier(emailVerificationRepo)

	// user
	register := user.NewRegister(codeVerifier, userRepo, hasher, publisher)
	userHandler := http.NewUserHandler(register)

	// auth
	login := aAuth.NewLogin(userRepo, hasher, jwtManager, refreshTokenRepo, refreshTokenManager)
	emailOTP := aAuth.NewRequestEmailOTP(emailVerificationRepo, otp, publisher)
	authHandler := http.NewAuthHandler(login, emailOTP)

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
