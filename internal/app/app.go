package app

import (
	"log/slog"

	httpapp "github.com/fumkaa/go-sub-aggregator-api/internal/app/http"
	"github.com/fumkaa/go-sub-aggregator-api/internal/config"
	httpserver "github.com/fumkaa/go-sub-aggregator-api/internal/transport/http"
)

type App struct {
	HttpServer *httpapp.App
}

func MustNew(log *slog.Logger, subService httpserver.SubscriptionService, cfg *config.Config) *App {
	return &App{
		HttpServer: httpapp.New(log, subService, &cfg.HttpServer),
	}
}
