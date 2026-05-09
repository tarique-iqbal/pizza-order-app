package container

import (
	"os"

	authapp "identity-service/internal/application/auth"
	"identity-service/internal/application/user"
	authinfra "identity-service/internal/infrastructure/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/redis"
	"identity-service/internal/infrastructure/security"
	"identity-service/internal/interfaces/http"
	"identity-service/internal/interfaces/http/middlewares"
)

type APIContainer struct {
	*Shared
	UserHandler *http.UserHandler
	AuthHandler *http.AuthHandler
	Middleware  *middlewares.Middleware
}

func NewAPIContainer() (*APIContainer, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	cfg := redis.Config{
		Addr: os.Getenv("REDIS_ADDR"),
	}

	base, err := NewShared()
	if err != nil {
		return nil, err
	}

	rc, err := redis.InitRedis(cfg)
	if err != nil {
		return nil, err
	}

	hasher := security.NewPasswordHasher()
	jwtManager := security.NewJWTManager(jwtSecret)
	refreshTokenManager := security.NewRefreshTokenManager()
	middleware := middlewares.NewMiddleware(jwtManager)
	otp := security.NewOTPGenerator()

	refreshTokenRepo := persistence.NewRefreshTokenRepository(rc)
	emailVerificationRepo := persistence.NewEmailVerificationRepository(base.DB)
	userRepo := persistence.NewUserRepository(base.DB)

	codeVerifier := authinfra.NewEmailVerifier(emailVerificationRepo)

	// user
	registerCustomer := user.NewRegisterCustomer(codeVerifier, userRepo, hasher, base.Publisher)
	registerOwner := user.NewRegisterOwner(base.DB, codeVerifier, hasher, userRepo, base.OutboxRepo, base.Publisher)
	findByID := user.NewFindByID(userRepo)
	userHandler := http.NewUserHandler(registerCustomer, registerOwner, findByID)

	// auth
	login := authapp.NewLogin(userRepo, hasher, jwtManager, refreshTokenRepo, refreshTokenManager)
	emailOTP := authapp.NewRequestEmailOTP(emailVerificationRepo, otp, base.Publisher)
	refreshToken := authapp.NewRefreshToken(jwtManager, refreshTokenRepo, refreshTokenManager)
	logout := authapp.NewLogout(refreshTokenRepo, refreshTokenManager)
	authHandler := http.NewAuthHandler(login, emailOTP, refreshToken, logout)

	return &APIContainer{
		Shared:      base,
		UserHandler: userHandler,
		AuthHandler: authHandler,
		Middleware:  middleware,
	}, nil
}
