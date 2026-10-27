package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sahilrana27582/payment-system/internal/config"
)

func RunMigrations(cfg *config.Config) error {
	databaseURL := cfg.DatabaseURL()

	m, err := migrate.New(
		"file://internal/db/migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
