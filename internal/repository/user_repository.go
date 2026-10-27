package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

var ErrDuplicateEmail = errors.New("email already exists")

type UserRepository struct {
	Db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{Db: db}
}

func (repo *UserRepository) CreateUser(ctx context.Context, username, password, email string) (string, error) {
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
