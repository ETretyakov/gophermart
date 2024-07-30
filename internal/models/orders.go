package models

import (
	"gophermart/internal/types"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         string `json:"order_id"    db:"order_id"`
	UserID     string `json:"user_id"     db:"user_id"`
	Number     uint64 `json:"number"      db:"number"`
	Status     string `json:"status"      db:"status"`
	Accrual    int64  `json:"accrual"     db:"accrual"`
	UploadedAt string `json:"uploaded_at" db:"uploaded_at"`
}

func NewOrder(userID string, number uint64) *Order {
	return &Order{
		ID:         uuid.NewString(),
		UserID:     userID,
		Number:     number,
		Status:     types.OrderNew.String(),
		Accrual:    0,
		UploadedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

type AccrueRecord struct {
	UserID string
	Number uint64
	Amount float64
}
