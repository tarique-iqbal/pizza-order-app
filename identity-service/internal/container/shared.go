package container

import (
	"os"

	"gorm.io/gorm"

	"identity-service/internal/domain/outbox"
	"identity-service/internal/infrastructure/db"
	"identity-service/internal/infrastructure/messaging"
	"identity-service/internal/infrastructure/persistence"
)

type Shared struct {
	DB         *gorm.DB
	OutboxRepo outbox.OutboxRepository
	Publisher  *messaging.RabbitMQPublisher
}

func NewShared() (*Shared, error) {
	amqpURL := os.Getenv("RABBITMQ_URL")

	database, err := db.InitDB()
	if err != nil {
		return nil, err
	}

	outboxRepo := persistence.NewOutboxRepository(database)
	publisher := messaging.NewRabbitMQPublisher(amqpURL)

	return &Shared{
		DB:         database,
		OutboxRepo: outboxRepo,
		Publisher:  publisher,
	}, nil
}

func (c *Shared) Close() {
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
