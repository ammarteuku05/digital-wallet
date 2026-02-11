package di

import (
	"digital-wallet/configs"
	"log/slog"
	"os"
)

func newLogger(cfg *configs.Config) *slog.Logger {
	var level slog.Level

	switch cfg.Logger.Level {
	case 0:
		level = slog.LevelDebug
	case 1:
		level = slog.LevelInfo
	case 2:
		level = slog.LevelWarn
	case 3:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
