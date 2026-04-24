package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fumkaa/go-sub-aggregator-api/internal/config"
	httpserver "github.com/fumkaa/go-sub-aggregator-api/internal/transport/http"
)

type App struct {
	server *http.Server
	log    *slog.Logger
	port   int
}

func New(log *slog.Logger, subService httpserver.SubscriptionService, config *config.HttpServerConfig) *App {
	api := httpserver.NewServerAPI(subService, log)
	router := httpserver.RegisterRoutes(api)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	return &App{
		server: server,
		log:    log,
		port:   config.Port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "httpapp.App.Run()"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	log.Info("http server is running...", slog.String("address", a.server.Addr))

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) {
	const op = "httpapp.App.Stop()"

	log := a.log.With(slog.String("op", op))

	log.Info("stopping http server...", slog.Int("port", a.port))

	a.server.Shutdown(ctx)

	log.Info("http server stopped.")
}
