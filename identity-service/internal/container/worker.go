package container

import (
	"log/slog"

	"identity-service/internal/application/outbox"
)

type WorkerContainer struct {
	*Shared
	Worker *outbox.Worker
}

func NewWorkerContainer(logger *slog.Logger) (*WorkerContainer, error) {
	base, err := NewShared()
	if err != nil {
		return nil, err
	}

	relayer := outbox.NewRelay(base.Publisher)
	worker := outbox.NewWorker(base.OutboxRepo, relayer, outbox.DefaultConfig(), logger)

	return &WorkerContainer{
		Shared: base,
		Worker: worker,
	}, nil
}
