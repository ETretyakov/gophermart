package handlers

import (
	"encoding/json"
	"gophermart/internal/controllers"
	"gophermart/internal/crypto"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"gophermart/internal/models"
	"gophermart/internal/repository"
	"gophermart/internal/validators"
	"net/http"

	"github.com/pkg/errors"
)

type BalanceHandlers struct {
	validator  validators.BalanceValidator
	controller controllers.BalanceController
	logger     log.HTTPLogger
}

func NewBalanceHandlers(repos *repository.Repos) *BalanceHandlers {
	return &BalanceHandlers{
		validator:  validators.NewBalanceValidator(),
		controller: controllers.NewBalanceController(repos),
		logger:     log.NewHTTPLogger("BalanceHandlers"),
	}
}

func (h *BalanceHandlers) GetForUser(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	balance, err := h.controller.GetForUser(ctx, userClaims.Subject)
	if err != nil {
		h.logger.Error(r, "failed to get user balance", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balanceOut := models.BalanceRead{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}
	if err := json.NewEncoder(w).Encode(&balanceOut); err != nil {
		h.logger.Error(r, "failed to encode response json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
