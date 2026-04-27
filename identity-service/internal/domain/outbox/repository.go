package outbox

import (
	"context"

	"gorm.io/gorm"
)

type OutboxRepository interface {
	WithTx(tx *gorm.DB) OutboxRepository
	Create(ctx context.Context, outbox *OutboxEvent) error
	FetchAndMarkProcessing(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkProcessed(ctx context.Context, id int64) error
	ReleaseForRetry(ctx context.Context, id int64, errMsg string) error
}
