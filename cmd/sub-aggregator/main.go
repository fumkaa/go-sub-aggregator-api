package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fumkaa/go-sub-aggregator-api/internal/app"
	"github.com/fumkaa/go-sub-aggregator-api/internal/config"
)

const (
	envLocal = "local"
	envProd  = "production"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	logger.Info("starting sub-aggregator api",
		slog.String("env", cfg.Env),
		// remove in production
		slog.Any("config", cfg),
		// ---
		slog.Int("port", cfg.HttpServer.Port),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	application := app.MustNew(logger, nil, cfg)

	go application.HttpServer.MustRun()

	// Graceful shutdown
	<-ctx.Done()

	logger.Info("stopping sub-aggregator api...")

	application.HttpServer.Stop(ctx)

	logger.Info("sub-aggregator api stopped.")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	return logger
}
