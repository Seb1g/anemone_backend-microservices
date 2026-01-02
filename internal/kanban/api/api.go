package api

import (
	"anemone_backend-microservices/internal/kanban/model"
	"anemone_backend-microservices/internal/kanban/repository"
	"anemone_backend-microservices/internal/kanban/services"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if errors.Is(err, repository.ErrBoardNotFound) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Board not found"})
		return
	}

	var reqErr *json.SyntaxError
	if errors.As(err, &reqErr) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload"})
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
}

type Handler struct {
	Service   *services.Service
	BoardRepo BoardRepoInterface
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) Routes(r *mux.Router, JWTsecret string) {
	boardRepo := h.BoardRepo

	api := r.PathPrefix("/api/v1/trello").Subrouter()
	api.Use(AuthMiddleware(JWTsecret))

	// --- BOARDS ---
	api.HandleFunc("/create_board", h.createBoard).Methods("POST")
	api.HandleFunc("/get_all_user_boards", h.getAllUserBoards).Methods("GET")
	boardRouter := api.PathPrefix("/board/{boardID}").Subrouter()
	boardRouter.Use(func(next http.Handler) http.Handler {
		return IsBoardOwner_Path(boardRepo, next)
	})

	boardRouter.Handle("", http.HandlerFunc(h.getOneUserBoard)).Methods("GET")
	boardRouter.Handle("", http.HandlerFunc(h.deleteBoard)).Methods("DELETE")
	boardRouter.Handle("", http.HandlerFunc(h.renameBoard)).Methods("PUT")
	boardRouter.Handle("", http.HandlerFunc(h.updateBoard)).Methods("POST")

	// --- COLUMNS ---
	boardRouter.HandleFunc("/column", h.createColumn).Methods("POST")

	columnRouter := boardRouter.PathPrefix("/column/{columnID}").Subrouter()
	columnRouter.Handle("", http.HandlerFunc(h.deleteColumn)).Methods("DELETE")
	columnRouter.Handle("", http.HandlerFunc(h.renameColumn)).Methods("PUT")

	// --- CARDS ---
	columnContext := api.PathPrefix("/column/{columnID}").Subrouter()
	columnContext.Use(func(next http.Handler) http.Handler {
		return IsBoardOwner_ColumnPath(boardRepo, next)
	})

	columnContext.HandleFunc("/card", h.createCard).Methods("POST")

	cardRouter := columnContext.PathPrefix("/card/{cardID}").Subrouter()
	cardRouter.Handle("", http.HandlerFunc(h.deleteCard)).Methods("DELETE")
	cardRouter.Handle("", http.HandlerFunc(h.renameCard)).Methods("PUT")
}

/* ====================== BOARD ====================== */

func (h *Handler) createBoard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}
	defer r.Body.Close()

	userID := UserID(r.Context())

	boardData, err := h.Service.CreateBoard(r.Context(), req.Title, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boardData)
}

func (h *Handler) getOneUserBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID := vars["boardID"]

	oneUserBoard, err := h.Service.GetOneUserBoard(r.Context(), boardID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oneUserBoard)
}

func (h *Handler) getAllUserBoards(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r.Context())

	allUserBoards, err := h.Service.GetAllUserBoards(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allUserBoards)
}

func (h *Handler) deleteBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID := vars["boardID"]

	err := h.Service.DeleteBoard(r.Context(), boardID)
	if err != nil {
		if errors.Is(err, repository.ErrBoardNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Board not found"})
			return
		}
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Board deleted successfully"})
}

func (h *Handler) renameBoard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NewName string `json:"new_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}
	defer r.Body.Close()

	vars := mux.Vars(r)
	boardID := vars["boardID"]

	board, err := h.Service.RenameBoard(r.Context(), boardID, req.NewName)
	if err != nil {
		if errors.Is(err, repository.ErrBoardNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Board not found"})
			return
		}
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

func (h *Handler) updateBoard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BoardData []*trello_model.Column `json:"board_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}
	defer r.Body.Close()

	vars := mux.Vars(r)
	boardID := vars["boardID"]

	userID := UserID(r.Context())

	err := h.Service.UpdateBoard(r.Context(), boardID, userID, req.BoardData)
	if err != nil {
		if errors.Is(err, repository.ErrBoardNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Board not found"})
			return
		}
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"message": "Board updated successfully", "success": true})
}

/* ====================== COLUMN ====================== */

func (h *Handler) createColumn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID := vars["boardID"]

	var req struct {
		ColumnTitle string `json:"column_title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}
	defer r.Body.Close()

	columnData, err := h.Service.CreateColumn(r.Context(), boardID, req.ColumnTitle)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(columnData)
}

func (h *Handler) deleteColumn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID := vars["boardID"]
	columnID := vars["columnID"]

	err := h.Service.DeleteColumn(r.Context(), boardID, columnID)
	if err != nil {
		if errors.Is(err, repository.ErrColumnNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Column not found"})
			return
		}
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Column deleted successfully"})
}

func (h *Handler) renameColumn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID := vars["boardID"]
	columnID := vars["columnID"]

	var req struct {
		NewName string `json:"new_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}
	defer r.Body.Close()

	column, err := h.Service.RenameColumn(r.Context(), boardID, columnID, req.NewName)
	if err != nil {
		if errors.Is(err, repository.ErrColumnNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Column not found"})
			return
		}
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(column)
}

/* ====================== CARD ====================== */

func (h *Handler) createCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	columnID := vars["columnID"]

	var req struct {
		CardTitle string `json:"card_title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}
	defer r.Body.Close()

	cardData, err := h.Service.CreateCard(r.Context(), columnID, req.CardTitle)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cardData)
}

func (h *Handler) deleteCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	columnID := vars["columnID"]
	cardID := vars["cardID"]

	err := h.Service.DeleteCard(r.Context(), columnID, cardID)
	if err != nil {
		if errors.Is(err, repository.ErrCardNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Card not found"})
			return
		}
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Card deleted successfully"})
}

func (h *Handler) renameCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	columnID := vars["columnID"]
	cardID := vars["cardID"]

	var req struct {
		NewName string `json:"new_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}
	defer r.Body.Close()

	card, err := h.Service.RenameCard(r.Context(), columnID, cardID, req.NewName)
	if err != nil {
		if errors.Is(err, repository.ErrCardNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Card not found"})
			return
		}
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"message": "Card success renamed", "card": card})
}
