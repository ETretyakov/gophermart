package models

import (
	"time"

	"github.com/google/uuid"
)

type Balance struct {
	ID        string  `json:"balance_id" db:"balance_id"`
	UserID    string  `json:"user_id"    db:"user_id"`
	Current   float64 `json:"current"    db:"current"`
	Withdrawn float64 `json:"withdrawn"  db:"withdrawn"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
}

func NewBalance(userID string) *Balance {
	return &Balance{
		ID:        uuid.NewString(),
		UserID:    userID,
		Current:   0,
		Withdrawn: 0,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

type BalanceRead struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
