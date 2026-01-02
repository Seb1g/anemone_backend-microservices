package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"anemone_backend-microservices/internal/pkg"

	"github.com/gorilla/mux"
)

type ctxKey string

const EntityIDKey ctxKey = "entity_id"
const userIDKey ctxKey = "user_id"

func UserID(ctx context.Context) int64 {
	return ctx.Value(userIDKey).(int64)
}

type OwnerChecker interface {
	IsNoteOwner(noteID int, userID int64) (bool, error)
	IsFolderOwner(folderID int, userID int64) (bool, error)
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

func NoteOwner(repo OwnerChecker) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := UserID(r.Context())

			idStr := mux.Vars(r)["id"]
			noteID, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
				return
			}

			ok, err := repo.IsNoteOwner(noteID, userID)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), EntityIDKey, noteID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FolderOwner(repo OwnerChecker) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := UserID(r.Context())

			idStr := mux.Vars(r)["id"]
			folderID, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
				return
			}

			ok, err := repo.IsFolderOwner(folderID, userID)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), EntityIDKey, folderID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
