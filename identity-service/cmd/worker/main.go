package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"identity-service/internal/container"
	"identity-service/internal/logger"
)

func main() {
	l := logger.New()
	slog.SetDefault(l)

	c, err := container.NewWorkerContainer(l)
	if err != nil {
		slog.Error("failed to initialize worker container", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	var wg sync.WaitGroup

	slog.Info("starting outbox worker...")

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Worker.Start(ctx)
	}()

	<-ctx.Done()

	slog.Info("shutdown signal received")

	wg.Wait()

	slog.Info("worker stopped gracefully")
}
