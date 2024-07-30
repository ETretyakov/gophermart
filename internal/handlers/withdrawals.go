package handlers

import (
	"encoding/json"
	"gophermart/internal/controllers"
	"gophermart/internal/crypto"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"gophermart/internal/repository"
	"gophermart/internal/validators"
	"net/http"

	"github.com/pkg/errors"
)

type WithdrawalsHandlers struct {
	validator  validators.WithdrawalsValidator
	controller controllers.WithdrawalController
	logger     log.HTTPLogger
}

func NewWithdrawalsHandlers(repos *repository.Repos) *WithdrawalsHandlers {
	return &WithdrawalsHandlers{
		validator:  validators.NewWithdrawalsValidator(),
		controller: controllers.NewWithdrawalsController(repos),
		logger:     log.NewHTTPLogger("WithdrawalsHandlers"),
	}
}

func (h *WithdrawalsHandlers) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userClaims, err := crypto.JWT.AuthToken(r.Header)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrNotAuthorised):
			h.logger.Debug(r, "failed to auth user: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
		default:
			h.logger.Error(r, "failed to auth user", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	withdrawalIn, err := h.validator.ValidateOrderCreate(
		userClaims.Subject,
		r.Body,
	)

	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrWrongOrderNumber):
			h.logger.Debug(r, "order already accepted: %s", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
		default:
			h.logger.Debug(r, "failed to validate order body: %s", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	withdrawal, err := h.controller.Create(ctx, withdrawalIn)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrBalanceIsNegative):
			h.logger.Debug(r, "balance is negative: %s", err)
			w.WriteHeader(http.StatusPaymentRequired)
		default:
			h.logger.Error(r, "failed to auth user", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&withdrawal); err != nil {
		h.logger.Error(r, "failed to encode response json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *WithdrawalsHandlers) UserWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userClaims, err := crypto.JWT.AuthToken(r.Header)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrNotAuthorised):
			h.logger.Debug(r, "failed to auth user: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
		default:
			h.logger.Error(r, "failed to auth user", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	withdrawals, err := h.controller.UserWithdrawals(ctx, userClaims.Subject)
	if err != nil {
		h.logger.Error(r, "failed to get orders", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(*withdrawals) == 0 {
		h.logger.Debug(r, "no withdrawals found: %s", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&withdrawals); err != nil {
		h.logger.Error(r, "failed to encode response json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
