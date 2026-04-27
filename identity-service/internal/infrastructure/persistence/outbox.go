package persistence

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"identity-service/internal/domain/outbox"
)

type outboxRepo struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) outbox.OutboxRepository {
	return &outboxRepo{db: db}
}

func (repo *outboxRepo) WithTx(tx *gorm.DB) outbox.OutboxRepository {
	return &outboxRepo{db: tx}
}

func (repo *outboxRepo) Create(ctx context.Context, out *outbox.OutboxEvent) error {
	return repo.db.WithContext(ctx).Create(out).Error
}

func (repo *outboxRepo) FetchAndMarkProcessing(ctx context.Context, limit int) ([]outbox.OutboxEvent, error) {
	var events []outbox.OutboxEvent

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		table := (outbox.OutboxEvent{}).TableName()

		query := fmt.Sprintf(`
			UPDATE %s
			SET status = 'processing',
			    locked_until = NOW() + INTERVAL '30 seconds',
			    attempts = attempts + 1
			WHERE id IN (
			    SELECT id
			    FROM %s
			    WHERE status = 'pending'
			    ORDER BY created_at
			    LIMIT ?
			    FOR UPDATE SKIP LOCKED
			)
			RETURNING *
		`, table, table)

		return tx.Raw(query, limit).Scan(&events).Error
	})

	return events, err
}

func (repo *outboxRepo) MarkProcessed(ctx context.Context, id int64) error {
	return repo.db.WithContext(ctx).
		Model(&outbox.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "processed",
			"processed_at": time.Now().UTC(),
			"locked_until": nil,
		}).Error
}

func (repo *outboxRepo) ReleaseForRetry(ctx context.Context, id int64, errMsg string) error {
	return repo.db.WithContext(ctx).
		Model(&outbox.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "pending",
			"last_error":   errMsg,
			"locked_until": nil,
		}).Error
}
