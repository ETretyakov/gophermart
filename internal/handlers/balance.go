package handlers

import (
	"encoding/json"
	"gophermart/internal/controllers"
	"gophermart/internal/log"
	"gophermart/internal/middlewares"
	"gophermart/internal/models"
	"gophermart/internal/repository"
	"gophermart/internal/validators"
	"net/http"
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

	rawUserID := ctx.Value(middlewares.UserIDKey)
	userID, ok := rawUserID.(string)
	if !ok || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	balance, err := h.controller.GetForUser(ctx, userID)
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

	w.WriteHeader(http.StatusOK)
}
