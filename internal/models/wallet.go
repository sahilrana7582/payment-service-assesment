package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Wallet struct {
	ID        string         `json:"id" db:"id"`
	UserID    string         `json:"user_id" db:"user_id"`
	Balance   pgtype.Numeric `json:"balance" db:"balance"`
	Active    bool           `json:"active" db:"active"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}
