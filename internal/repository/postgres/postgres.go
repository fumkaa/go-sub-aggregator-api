package postgres

import (
	"context"
	"log/slog"

	"github.com/fumkaa/go-sub-aggregator-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	dbpool *pgxpool.Pool
	log    *slog.Logger
}

func MustNew(ctx context.Context, cfg *config.StorageConfig, log *slog.Logger) *Storage {
	const op = "postgres.MustNew()"
	tlog := log.With(slog.String("op", op))

	tlog.Info("creating postgres connection...")

	config, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		panic(err)
	}

	config.MaxConns = cfg.MaxConns
	config.MinConns = cfg.MinConns
	config.MaxConnIdleTime = cfg.MaxConnIdleTime
	config.HealthCheckPeriod = cfg.HealthCheckPeriod

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}

	tlog.Info("pinging postgres connection...")
	if err := dbpool.Ping(ctx); err != nil {
		panic(err)
	}

	tlog.Info("postgres connection created")
	return &Storage{
		dbpool: dbpool,
		log:    log,
	}
}

func (s *Storage) Close() {
	const op = "postgres.Close()"
	log := s.log.With(slog.String("op", op))

	log.Info("closing postgres connection...")
	s.dbpool.Close()
	log.Info("postgres connection closed.")
}
