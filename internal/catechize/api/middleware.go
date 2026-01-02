package api

import (
	"context"
	"net/http"
	"strings"

	"anemone_backend-microservices/internal/pkg"
)

type ctxKey string
const userIDKey ctxKey = "user_id"

func UserID(ctx context.Context) int64 {
	return ctx.Value(userIDKey).(int64)
}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			parts := strings.SplitN(h, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userID, err := auth_utils.ValidateToken(parts[1], secret)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}