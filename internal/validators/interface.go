package validators

import (
	"gophermart/internal/models"
	"io"
	"net/http"
)

type AuthValidator interface {
	ValidateUserRegister(body io.ReadCloser) (*models.AuthUser, error)
	ValidateUserLogin(body io.ReadCloser) (*models.UserLogin, error)
}

type OrderValidator interface {
	ValidateOrderCreate(userID string, body io.ReadCloser) (*models.Order, error)
	ValidateOrderFromPath(r *http.Request) (string, error)
}

type BalanceValidator interface{}

type WithdrawalsValidator interface {
	ValidateOrderCreate(userID string, body io.ReadCloser) (*models.Withdrawal, error)
}
