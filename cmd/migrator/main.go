package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		storagePath     string
		migrationsPath  string
		migrationsTable string
		down            bool
	)

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migration table")
	flag.BoolVar(&down, "down", false, "run down migrations")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is required")
	}

	if migrationsPath == "" {
		panic("migrations path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("%s?x-migrations-table=%s&sslmode=disable", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}
	defer m.Close()

	if down {
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("no migrations to rollback")
				return
			}
			panic(err)
		}
		log.Println("migrations rolled back successfully")
		return
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	log.Println("migrations applied successfully")
}
