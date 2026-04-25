package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fumkaa/go-sub-aggregator-api/internal/app"
	"github.com/fumkaa/go-sub-aggregator-api/internal/config"
	"github.com/fumkaa/go-sub-aggregator-api/internal/lib/logger/sl"
	"github.com/fumkaa/go-sub-aggregator-api/internal/repository/postgres"
	"github.com/fumkaa/go-sub-aggregator-api/internal/services/subscription"
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
		slog.Int("port", cfg.HttpServer.Port),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	storage := postgres.MustNew(ctx,
		&cfg.Storage,
		logger,
	)

	subscriptionService := subscription.NewSubscriptionManager(logger, storage, storage, storage)
	application := app.MustNew(logger, subscriptionService, cfg)

	go application.HttpServer.MustRun()

	// Graceful shutdown
	<-ctx.Done()

	logger.Info("stopping sub-aggregator api...")

	application.HttpServer.Stop(ctx)
	storage.Close()

	logger.Info("sub-aggregator api stopped.")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		logger = slog.New(&sl.ContextHandler{Handler: baseHandler})
	case envProd:
		baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger = slog.New(&sl.ContextHandler{Handler: baseHandler})
	}
	slog.SetDefault(logger)
	return logger
}
