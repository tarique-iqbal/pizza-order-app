package fixtures

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"identity-service/internal/domain/outbox"
	"identity-service/tests/testutil"
)

func LoadOutboxEventFixtures(db *gorm.DB) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range 5 {
		restaurantID, payload, err := restaurantPayload(r)
		if err != nil {
			return err
		}

		event := outbox.NewOutboxEvent(
			restaurantID,
			outbox.EventRestaurantInitiated,
			payload,
		)

		if err := db.Create(&event).Error; err != nil {
			return err
		}
	}

	return nil
}

func restaurantPayload(r *rand.Rand) (restaurantID uuid.UUID, payload []byte, err error) {
	restaurantID = testutil.MustNewID()

	payloadMap := map[string]any{
		"restaurant_id": restaurantID,
		"owner_id":      testutil.MustNewID(),
		"business_name": randomString(r, "Restaurant"),
		"vat_number":    randomVAT(r),
	}

	payload, err = json.Marshal(payloadMap)
	if err != nil {
		return restaurantID, nil, err
	}

	return restaurantID, payload, nil
}

func randomString(r *rand.Rand, prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, r.Intn(1_000_000))
}

func randomVAT(r *rand.Rand) string {
	return fmt.Sprintf("DE%09d", r.Intn(1_000_000_000))
}
