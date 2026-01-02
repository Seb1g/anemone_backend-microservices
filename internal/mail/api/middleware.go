package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"anemone_backend-microservices/internal/mail/repository"
	"anemone_backend-microservices/internal/pkg"

	"github.com/gorilla/mux"
)

type ctxKey string

const (
	UserIDKey    ctxKey = "user_id"
	AddressIDKey ctxKey = "address_id"
)

func GetUserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(UserIDKey).(int64)
	return id, ok
}

func Auth(secret string) mux.MiddlewareFunc {
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

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CheckAddressOwner(repo *repository.Repository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			idStr := mux.Vars(r)["id"]
			addressID, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "bad address id", http.StatusBadRequest)
				return
			}

			ok, err = repo.CheckAddressOwner(addressID, userID)
			if err != nil || !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), AddressIDKey, addressID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
