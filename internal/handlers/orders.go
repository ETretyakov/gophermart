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

type OrdersHandlers struct {
	validator  validators.OrderValidator
	controller controllers.OrdersController
	logger     log.HTTPLogger
}

func NewOrdersHandlers(repos *repository.Repos) *OrdersHandlers {
	return &OrdersHandlers{
		validator:  validators.NewOrdersValidator(),
		controller: controllers.NewOrdersController(repos),
		logger:     log.NewHTTPLogger("OrdersHandlers"),
	}
}

func (h *OrdersHandlers) Create(w http.ResponseWriter, r *http.Request) {
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

	orderIn, err := h.validator.ValidateOrderCreate(userClaims.Subject, r.Body)
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

	_, err = h.controller.Create(ctx, orderIn)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrOrderAlreadyAccepted):
			h.logger.Debug(r, "order already accepted: %s", err)
			w.WriteHeader(http.StatusOK)
		case errors.Is(err, exceptions.ErrOrderAlreadyRegistered):
			h.logger.Debug(r, "order already registered: %s", err)
			w.WriteHeader(http.StatusConflict)
		default:
			h.logger.Error(r, "failed to auth user", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *OrdersHandlers) UserOrders(w http.ResponseWriter, r *http.Request) {
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

	orders, err := h.controller.UserOrders(ctx, userClaims.Subject)
	if err != nil {
		h.logger.Error(r, "failed to get orders", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(*orders) == 0 {
		h.logger.Debug(r, "no orders found: %s", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&orders); err != nil {
		h.logger.Error(r, "failed to encode response json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *OrdersHandlers) UserOrderByNumber(w http.ResponseWriter, r *http.Request) {
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

	orderNumber, err := h.validator.ValidateOrderFromPath(r)
	if err != nil {
		h.logger.Debug(r, "failed to parse order: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	order, err := h.controller.GetUserOrderByNumber(ctx, userClaims.Subject, *orderNumber)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrOrderNotFound):
			h.logger.Debug(r, "failed to find order: %s", err)
			w.WriteHeader(http.StatusNoContent)
		default:
			h.logger.Error(r, "failed to get orders", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&order); err != nil {
		h.logger.Error(r, "failed to encode response json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
