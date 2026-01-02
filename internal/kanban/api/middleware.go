package api

import (
	"anemone_backend-microservices/internal/kanban/repository"
	"context"
	"errors"
	"net/http"
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

type BoardRepoInterface interface {
	GetBoardOwnerID(ctx context.Context, boardID string) (int64, error)
	GetBoardOwnerIDByColumnID(ctx context.Context, columnID string) (int64, error)
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

func IsBoardOwner_Query(boardRepo BoardRepoInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		boardID := r.URL.Query().Get("boardId")
		if boardID == "" {
			http.Error(w, "Board ID is missing in query parameters", http.StatusBadRequest)
			return
		}

		userID := UserID(r.Context())

		ownerID, err := boardRepo.GetBoardOwnerID(r.Context(), boardID)
		if err != nil {
			if errors.Is(err, repository.ErrBoardNotFound) {
				http.Error(w, "Board not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Error checking board ownership", http.StatusInternalServerError)
			return
		}

		if int64(ownerID) != userID {
			http.Error(w, "Access Forbidden: Not the board owner", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IsBoardOwner_Path(boardRepo BoardRepoInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		boardID, ok := vars["boardID"]
		if !ok || boardID == "" {
			http.Error(w, "Board ID is missing in URL path", http.StatusBadRequest)
			return
		}

		userID := UserID(r.Context())

		ownerID, err := boardRepo.GetBoardOwnerID(r.Context(), boardID)
		if err != nil {
			if errors.Is(err, repository.ErrBoardNotFound) {
				http.Error(w, "Board not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Error checking board ownership", http.StatusInternalServerError)
			return
		}

		if int64(ownerID) != userID {
			http.Error(w, "Access Forbidden: Not the board owner", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IsBoardOwner_ColumnPath(boardRepo BoardRepoInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		columnID, ok := vars["columnID"]
		if !ok || columnID == "" {
			http.Error(w, "Column ID is missing in URL path", http.StatusBadRequest)
			return
		}

		userID := UserID(r.Context())

		ownerID, err := boardRepo.GetBoardOwnerIDByColumnID(r.Context(), columnID)
		if err != nil {
			if errors.Is(err, repository.ErrBoardNotFound) {
				http.Error(w, "Column or Board not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Error checking board ownership", http.StatusInternalServerError)
			return
		}

		if ownerID != userID {
			http.Error(w, "Access Forbidden: Not the board owner", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
