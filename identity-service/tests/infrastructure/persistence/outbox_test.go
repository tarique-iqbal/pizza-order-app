package persistence_test

import (
	"context"
	"encoding/json"
	"identity-service/internal/domain/outbox"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/tests/infrastructure/db/fixtures"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOutboxRepo() outbox.OutboxRepository {
	ts := testStorage()
	truncateTables(ts.DB)

	if err := fixtures.LoadOutboxEventFixtures(ts.DB); err != nil {
		panic(err)
	}

	return persistence.NewOutboxRepository(ts.DB)
}

func TestOutboxRepository_Create(t *testing.T) {
	ts := testStorage()
	repo := setupOutboxRepo()

	restaurantID := generateUUID()
	payloadMap := map[string]interface{}{
		"restaurant_id": restaurantID,
		"owner_id":      generateUUID(),
		"business_name": "Domino's Pizza",
		"vat_number":    "DE123456121",
	}
	payload, _ := json.Marshal(payloadMap)

	event := outbox.NewOutboxEvent(
		restaurantID,
		outbox.EventRestaurantInitiated,
		payload,
	)

	err := repo.Create(context.Background(), &event)

	require.NoError(t, err)
	assert.NotZero(t, event.ID)

	var found outbox.OutboxEvent
	err = ts.DB.First(&found, event.ID).Error
	require.NoError(t, err)

	assert.Equal(t, outbox.StatusPending, found.Status)
	assert.NotZero(t, found.CreatedAt)
}

func TestOutboxRepository_FetchAndMarkProcessing(t *testing.T) {
	ts := testStorage()
	repo := setupOutboxRepo()

	events, err := repo.FetchAndMarkProcessing(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, events, 1)

	e := events[0]

	// runtime checks
	assert.Equal(t, outbox.StatusProcessing, e.Status)
	assert.NotNil(t, e.LockedUntil)
	assert.Equal(t, 1, e.Attempts)

	// DB verification (critical)
	var dbEvent outbox.OutboxEvent
	err = ts.DB.First(&dbEvent, e.ID).Error
	require.NoError(t, err)

	assert.Equal(t, outbox.StatusProcessing, dbEvent.Status)
	assert.NotNil(t, dbEvent.LockedUntil)
	assert.Equal(t, 1, dbEvent.Attempts)
}

func TestOutboxRepository_FetchAndMarkProcessing_Limit(t *testing.T) {
	repo := setupOutboxRepo()

	// insert extra events
	for range 3 {
		e := outbox.NewOutboxEvent(
			generateUUID(),
			outbox.EventRestaurantInitiated,
			[]byte(`{}`),
		)
		require.NoError(t, repo.Create(context.Background(), &e))
	}

	events, err := repo.FetchAndMarkProcessing(context.Background(), 2)

	require.NoError(t, err)
	assert.Len(t, events, 2)
}

func TestOutboxRepository_MarkProcessed(t *testing.T) {
	ts := testStorage()
	repo := setupOutboxRepo()

	var event outbox.OutboxEvent
	require.NoError(t, ts.DB.First(&event).Error)

	err := repo.MarkProcessed(context.Background(), event.ID)
	require.NoError(t, err)

	var updated outbox.OutboxEvent
	require.NoError(t, ts.DB.First(&updated, event.ID).Error)

	assert.Equal(t, outbox.StatusProcessed, updated.Status)
	assert.NotNil(t, updated.ProcessedAt)
	assert.Nil(t, updated.LockedUntil)
}

func TestOutboxRepository_ReleaseForRetry(t *testing.T) {
	ts := testStorage()
	repo := setupOutboxRepo()

	var event outbox.OutboxEvent
	require.NoError(t, ts.DB.First(&event).Error)

	// simulate processing state
	require.NoError(t, ts.DB.Model(&event).Updates(map[string]interface{}{
		"status":       outbox.StatusProcessing,
		"locked_until": time.Now().Add(30 * time.Second),
	}).Error)

	err := repo.ReleaseForRetry(context.Background(), event.ID, "temporary failure")
	require.NoError(t, err)

	var updated outbox.OutboxEvent
	require.NoError(t, ts.DB.First(&updated, event.ID).Error)

	require.NotNil(t, updated.LastError)
	assert.Equal(t, outbox.StatusPending, updated.Status)
	assert.Equal(t, "temporary failure", *updated.LastError)
	assert.Nil(t, updated.LockedUntil)
}

func TestOutboxRepository_FetchAndMarkProcessing_SkipLocked(t *testing.T) {
	ts := testStorage()
	repo := setupOutboxRepo()

	var event outbox.OutboxEvent
	require.NoError(t, ts.DB.First(&event).Error)

	// lock one event manually
	require.NoError(t, ts.DB.Model(&event).Updates(map[string]interface{}{
		"status":       outbox.StatusProcessing,
		"locked_until": time.Now().Add(30 * time.Second),
	}).Error)

	events, err := repo.FetchAndMarkProcessing(context.Background(), 10)

	require.NoError(t, err)

	// locked event should NOT be returned
	for _, e := range events {
		assert.NotEqual(t, event.ID, e.ID)
	}
}
