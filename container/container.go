package container

import (
	"log"
	"pizza-order-api/internal/application/user"
	"pizza-order-api/internal/infrastructure/db"
	"pizza-order-api/internal/infrastructure/persistence"
	"pizza-order-api/internal/interfaces/http"
	"pizza-order-api/internal/interfaces/http/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Container struct {
	UserHandler *http.UserHandler
	DB          *gorm.DB
}

func NewContainer() (*Container, error) {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("Could not connect to database:", err)
		return nil, err
	}

	userRepo := persistence.NewUserRepository(database)

	userUseCase := user.NewCreateUserUseCase(userRepo)

	userHandler := http.NewUserHandler(userUseCase)

	return &Container{
		UserHandler: userHandler,
		DB:          database,
	}, nil
}

func (c *Container) SetupRoutes(router *gin.Engine) {
	routes.SetupUserRoutes(router, c.UserHandler)
}

func (c *Container) Close() {
	db, _ := c.DB.DB()
	db.Close()
}
