package api

import (
	"anemone_backend-kanban/internal/domain/apperr"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (h *handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, apperr.ErrCardNotFound),
	     errors.Is(err, apperr.ErrColumnNotFound),
	     errors.Is(err, apperr.ErrBoardNotFound):
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return

	case errors.Is(err, io.EOF):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Request body cannot be empty"})
		return
	}

	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &unmarshalErr) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload"})
		return
	}

	fmt.Printf("[CRITICAL ERROR] Path: %s | Details: %v\n", r.URL.Path, err)

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
}