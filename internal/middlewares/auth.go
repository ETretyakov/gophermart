package middlewares

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/crypto"
	"gophermart/internal/exceptions"
	"gophermart/internal/log"
	"net/http"
)

type AuthKeyType string

var UserIDKey AuthKeyType = "user_id_key"

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userClaims, err := crypto.JWT.AuthToken(r.Header)
			if err != nil {
				switch {
				case errors.Is(err, exceptions.ErrNotAuthorised):
					log.Debug(ctx, fmt.Sprintf("failed to auth user: %s", err))
					w.WriteHeader(http.StatusUnauthorized)
				default:
					log.Error(ctx, "failed to auth user", err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			ctx = context.WithValue(ctx, UserIDKey, userClaims.Subject)

			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
