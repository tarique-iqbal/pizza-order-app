package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	if os.Getenv("APP_ENV") == "prod" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
