package container

import (
	"os"
	"restaurant-service/internal/infrastructure/db"

	"gorm.io/gorm"
)

type Shared struct {
	AMQPURL string
	DB      *gorm.DB
}

func NewShared() (*Shared, error) {
	database, err := db.InitDB()
	if err != nil {
		return nil, err
	}

	amqpURL := os.Getenv("RABBITMQ_URL")

	return &Shared{
		AMQPURL: amqpURL,
		DB:      database,
	}, nil
}
