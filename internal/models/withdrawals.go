package models

import (
	"github.com/google/uuid"
)

type Withdrawal struct {
	ID          string `json:"withdrawal_id" db:"withdrawal_id"`
	UserID      string `json:"user_id"       db:"user_id"`
	Order       string `json:"order"         db:"order"`
	Sum         int64  `json:"sum"           db:"sum"`
	ProcessedAt string `json:"processed_at"  db:"processed_at"`
}

func NewWithdrawal(userID string, order string, sum int64) *Withdrawal {
	return &Withdrawal{
		ID:     uuid.NewString(),
		UserID: userID,
		Order:  order,
		Sum:    sum,
	}
}

type WithdrawalCreate struct {
	Order string `json:"order"`
	Sum   int64  `json:"sum"`
}
