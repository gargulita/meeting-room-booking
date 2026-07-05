package middleware

import (
	"context"
	"net/http"
	"strings"

	"booking-service/internal/handlers"
	"booking-service/internal/utils"
)

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handlers.WriteAPIError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authorization header required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
				handlers.WriteAPIError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header format")
				return
			}

			claims, err := utils.ValidateJWT(parts[1], secret)
			if err != nil {
				handlers.WriteAPIError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), utils.CtxUserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, utils.CtxUserRoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
