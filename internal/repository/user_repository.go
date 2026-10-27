package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolationCode = "23505"

var ErrDuplicateEmail = errors.New("email already exists")

func (repo *Repository) CreateUser(ctx context.Context, username, password, email string) (string, error) {
	query := `
		INSERT INTO users (username, password, email)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var userID string
	if err := repo.Db.QueryRow(insertCtx, query, username, password, email).Scan(&userID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return "", ErrDuplicateEmail
		}
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func (repo *Repository) CreateUserWithWallet(ctx context.Context, username, password, email string) (string, error) {
	var userID string

	err := repo.WithTx(ctx, func(txRepo *Repository) error {
		var err error
		userID, err = txRepo.CreateUser(ctx, username, password, email)
		if err != nil {
			return err
		}

		return txRepo.CreateWallet(ctx, userID)
	})
	if err != nil {
		return "", err
	}

	return userID, nil
}
