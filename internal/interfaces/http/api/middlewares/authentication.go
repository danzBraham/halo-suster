package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	auth_error "github.com/danzBraham/halo-suster/internal/exceptions/auth"
	"github.com/danzBraham/halo-suster/internal/helpers"
)

type ContextKey string

var (
	ContextUserIDKey ContextKey = "userID"
	ContextRoleKey   ContextKey = "role"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helpers.ResponseJSON(w, http.StatusUnauthorized, &helpers.ResponseBody{
				Error:   "Unauthorized error",
				Message: "Missing Authorization header",
			})
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenString == "" {
			helpers.ResponseJSON(w, http.StatusUnauthorized, &helpers.ResponseBody{
				Error:   "Unauthorized error",
				Message: "Invalid Authorization header format",
			})
			return
		}

		credential, err := helpers.VerifyJWT(tokenString)
		if errors.Is(err, auth_error.ErrInvalidToken) {
			helpers.ResponseJSON(w, http.StatusUnauthorized, &helpers.ResponseBody{
				Error:   "Unauthorized error",
				Message: err.Error(),
			})
			return
		}
		if errors.Is(err, auth_error.ErrUnknownClaims) {
			helpers.ResponseJSON(w, http.StatusUnauthorized, &helpers.ResponseBody{
				Error:   "Unauthorized error",
				Message: err.Error(),
			})
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserIDKey, credential.UserId)
		ctx = context.WithValue(ctx, ContextRoleKey, credential.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
