package fixtures

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"identity-service/internal/domain/outbox"
)

func LoadOutboxEventFixtures(db *gorm.DB) error {
	restaurantID := uuid.MustParse("019da239-5f40-7212-8d20-eb3dde923e18")
	userID := uuid.MustParse("019dcf5d-d90e-7abc-8def-1234567890ab")

	payloadMap := map[string]interface{}{
		"restaurant_id": restaurantID,
		"owner_id":      userID,
		"business_name": "Test Restaurant",
		"vat_number":    "DE123456789",
	}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return err
	}

	events := []outbox.OutboxEvent{
		outbox.NewOutboxEvent(
			restaurantID,
			outbox.EventRestaurantInitiated,
			payload,
		),
	}

	for _, e := range events {
		e.CreatedAt = time.Now().UTC()
		e.Status = outbox.StatusPending
		e.Payload = datatypes.JSON(payload)

		if err := db.Create(&e).Error; err != nil {
			return err
		}
	}

	return nil
}
