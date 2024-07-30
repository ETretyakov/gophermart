package handlers

import (
	"encoding/json"
	"fmt"
	"gophermart/internal/controllers"
	"gophermart/internal/crypto"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"gophermart/internal/repository"
	"gophermart/internal/validators"
	"net/http"

	"github.com/pkg/errors"
)

type AuthHandlers struct {
	validator  validators.AuthValidator
	controller controllers.AuthController
	logger     log.HTTPLogger
}

func NewAuthHandlers(repos *repository.Repos) *AuthHandlers {
	return &AuthHandlers{
		validator:  validators.NewAuthValidator(),
		controller: controllers.NewAuthController(repos),
		logger:     log.NewHTTPLogger("AuthHandlers"),
	}
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authUser, err := h.validator.ValidateUserRegister(r.Body)
	if err != nil {
		h.logger.Debug(r, "failed to validate user register body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.controller.Register(ctx, authUser)
	if err != nil {
		if errors.Is(err, exceptions.ErrLoginAlreadyTaken) {
			h.logger.Debug(r, "failed to register user: %s", err)
			w.WriteHeader(http.StatusConflict)
		} else {
			h.logger.Error(r, "failed to register user", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(&user); err != nil {
		h.logger.Error(r, "failed to encode response json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := crypto.JWT.GetToken(authUser.ID)
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("failed to get token: %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token)

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	loginUser, err := h.validator.ValidateUserLogin(r.Body)
	if err != nil {
		h.logger.Debug(r, "failed to validate user login body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.controller.Login(ctx, loginUser)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrNotAuthorised):
			h.logger.Debug(r, "failed to login user: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
		case errors.Is(err, exceptions.ErrUserNotFound):
			h.logger.Debug(r, "failed to login user: %s", err)
			w.WriteHeader(http.StatusBadRequest)
		default:
			h.logger.Error(r, "failed to login user", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token)

	w.WriteHeader(http.StatusOK)
}
