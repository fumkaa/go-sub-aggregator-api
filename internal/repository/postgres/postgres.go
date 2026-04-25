package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	dbpool *pgxpool.Pool
	log    *slog.Logger
}

func MustNew(ctx context.Context, connString string, log *slog.Logger) *Storage {
	const op = "postgres.MustNew()"
	tlog := log.With(slog.String("op", op))

	tlog.Info("creating postgres connection...")

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnIdleTime = time.Minute * 30

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
