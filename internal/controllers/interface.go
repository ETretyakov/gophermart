package controllers

import (
	"context"
	"gophermart/internal/models"
)

type HealthController interface {
	SetReadiness(state bool)
	SetLiveness(state bool)
	ReadinessState() bool
	LivenessState() bool
	PingDB() bool
}

type AuthController interface {
	Register(ctx context.Context, schema *models.AuthUser) (*models.User, error)
	Login(ctx context.Context, schema *models.UserLogin) (string, error)
}

type OrdersController interface {
	Create(ctx context.Context, schema *models.Order) (*models.Order, error)
	UserOrders(ctx context.Context, userID string) (*[]models.Order, error)
	GetUserOrderByNumber(ctx context.Context, userID string, orderNumber uint64) (*models.Order, error)
}

type BalanceController interface {
	GetForUser(ctx context.Context, userID string) (*models.Balance, error)
}

type WithdrawalController interface {
	Create(ctx context.Context, schema *models.Withdrawal) (*models.Withdrawal, error)
	UserWithdrawals(ctx context.Context, userID string) (*[]models.Withdrawal, error)
}
