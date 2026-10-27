package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sahilrana27582/payment-system/internal/config"
)

var (
	once  sync.Once
	pool  *pgxpool.Pool
	dbErr error
)

func connect(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	databaseURL := cfg.DatabaseURL()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to the database")
	return pool, nil
}

func GetDBConn(cfg *config.Config) (*pgxpool.Pool, error) {
	once.Do(func() {
		ctx := context.Background()
		pool, dbErr = connect(ctx, cfg)
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return pool, nil
}

func PingDB(ctx context.Context, pool *pgxpool.Pool) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		return fmt.Errorf("database connection degraded: %w", err)
	}

	return nil
}
