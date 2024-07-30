package repository

import (
	"context"
	"errors"
	"gophermart/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repos struct {
	DB *sqlx.DB

	HealthRepo      HealthRepo
	AuthRepo        AuthRepo
	UsersRepo       UsersRepo
	BalanceRepo     BalanceRepo
	WithdrawalsRepo WithdrawalsRepo
	OrdersRepo      OrdersRepo
}

type HealthRepo interface {
	PingDB(ctx context.Context) bool
}

type AuthRepo interface {
	UserExists(ctx context.Context, email string) (bool, error)
	UserAuth(ctx context.Context, login string) (*models.AuthUser, error)
}

type UsersRepo interface {
	Create(ctx context.Context, model *models.AuthUser) (*models.User, error)
}

type BalanceRepo interface {
	GetOrCreateForUser(ctx context.Context, userID string) (*models.Balance, error)
	Get(ctx context.Context, balanceID string) (*models.Balance, error)
	Create(ctx context.Context, model *models.Balance) (*models.Balance, error)
	Update(ctx context.Context, model *models.Balance) (*models.Balance, error)
}

type WithdrawalsRepo interface {
	Create(ctx context.Context, model *models.Withdrawal) (*models.Withdrawal, error)
	UserWithdrawals(ctx context.Context, userID string) (*[]models.Withdrawal, error)
}

type OrdersRepo interface {
	Create(ctx context.Context, model *models.Order) (*models.Order, error)
	MarkAsProcessing(ctx context.Context, orderIDs []string) (bool, error)
	MarkAsInvalid(ctx context.Context, orderIDs []string) (bool, error)
	Accrue(ctx context.Context, record models.AccrueRecord) (bool, error)
	UserOrders(ctx context.Context, userID string) (*[]models.Order, error)
	GetByNumber(ctx context.Context, orderNumber string) (*models.Order, error)
	GetUserOrderByNumber(ctx context.Context, userID string, orderNumber uint64) (*models.Order, error)
}

func NewRepos(ctx context.Context, db *sqlx.DB) (*Repos, error) {
	if db != nil {
		repos := &Repos{DB: db}
		repos.HealthRepo = NewHealthRepo(repos)
		repos.AuthRepo = NewAuthRepo(repos)
		repos.UsersRepo = NewUsersRepo(repos)
		repos.BalanceRepo = NewBalanceRepo(repos)
		repos.WithdrawalsRepo = NewWithdrawalsRepo(repos)
		repos.OrdersRepo = NewOrdersRepoImpl(repos)
		return repos, nil
	} else {
		return nil, errors.New("database is not provided")
	}
}
