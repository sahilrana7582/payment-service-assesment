package repository

import (
	"context"
	"fmt"
	"time"
)

func (repo *Repository) CreateWallet(ctx context.Context, userID string) error {
	query := `
		INSERT INTO wallet (user_id)
		VALUES ($1);
	`

	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := repo.Db.Exec(insertCtx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to create wallet for user %s: %w", userID, err)
	}

	return nil
}
