package api

import (
	"anemone_backend-kanban/internal/domain/model"
	"encoding/json"
	"net/http"
)

/* ====================== CARD ====================== */

func (h *handler) createCard(w http.ResponseWriter, r *http.Request) {
	columnID := r.PathValue("column_id")

	if columnID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	var req struct {
		CardTitle string `json:"card_title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}
	defer r.Body.Close()

	cardData, err := h.CardService.CreateCard(r.Context(), columnID, req.CardTitle)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, cardData)
}

func (h *handler) deleteCard(w http.ResponseWriter, r *http.Request) {
	columnID := r.PathValue("column_id")
	cardID := r.PathValue("card_id")

	if columnID == "" || cardID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	err := h.CardService.DeleteCard(r.Context(), columnID, cardID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Card deleted successfully"})
}

func (h *handler) renameCard(w http.ResponseWriter, r *http.Request) {
	columnID := r.PathValue("column_id")
	cardID := r.PathValue("card_id")

	if columnID == "" || cardID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	var req struct {
		NewName string `json:"new_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}
	defer r.Body.Close()

	card, err := h.CardService.RenameCard(r.Context(), columnID, cardID, req.NewName)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]any{
		"message": "Card renamed successfully",
		"card":    card,
	})
}

func (h *handler) getAllCardsFromColumn(w http.ResponseWriter, r *http.Request) {
	columnID := r.PathValue("column_id")

	if columnID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	cards, err := h.CardService.GetAllCardsFromColumn(r.Context(), columnID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string][]*model.Card{
		"cards": cards,
	})
}

func (h *handler) moveCard(w http.ResponseWriter, r *http.Request) {
	columnID := r.PathValue("column_id")
	cardID := r.PathValue("card_id")

	if columnID == "" || cardID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	var req model.MoveCardParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}
	defer r.Body.Close()

	req.CardID = cardID
	req.FromColumnID = columnID

	err := h.CardService.MoveCard(r.Context(), req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "card moved successfully"})
}

/* ====================== COLUMN ====================== */

func (h *handler) createColumn(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")

	if boardID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	var req struct {
		ColumnTitle string `json:"column_title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}
	defer r.Body.Close()

	columnData, err := h.ColumnService.CreateColumn(r.Context(), boardID, req.ColumnTitle)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]*model.Column{
		"column": columnData,
	})
}

func (h *handler) deleteColumn(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	columnID := r.PathValue("column_id")

	if boardID == "" || columnID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	err := h.ColumnService.DeleteColumn(r.Context(), boardID, columnID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "Column deleted successfully",
	})
}

func (h *handler) renameColumn(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	columnID := r.PathValue("column_id")

	if boardID == "" || columnID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	var req struct {
		NewName string `json:"new_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}
	defer r.Body.Close()

	column, err := h.ColumnService.RenameColumn(r.Context(), boardID, columnID, req.NewName)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]*model.Column{
		"column": column,
	})
}

func (h *handler) getAllColumnsFromBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")

	if boardID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	columns, err := h.ColumnService.GetAllColumnsFromBoard(r.Context(), boardID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string][]*model.Column{
		"columns": columns,
	})
}

func (h *handler) moveColumn(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	columnID := r.PathValue("column_id")

	if boardID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	var req model.MoveColumnParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}
	defer r.Body.Close()

	req.BoardID = boardID
	req.ColumnID = columnID

	err := h.ColumnService.MoveColumn(r.Context(), req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "column move successfully",
	})
}

/* ====================== BOARD ====================== */

func (h *handler) createBoard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}
	defer r.Body.Close()

	userID := UserID(r.Context())

	boardData, err := h.BoardService.CreateBoard(r.Context(), req.Title, userID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]*model.Board{
		"board": boardData,
	})
}

func (h *handler) getAllUserBoards(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r.Context())

	allUserBoards, err := h.BoardService.GetAllUserBoards(r.Context(), userID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string][]*model.Board{
		"boards": allUserBoards,
	})
}

func (h *handler) deleteBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")

	if boardID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	err := h.BoardService.DeleteBoard(r.Context(), boardID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "board deleted successfully",
	})
}

func (h *handler) renameBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")

	if boardID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	var req struct {
		NewName string `json:"new_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err)
		return
	}
	defer r.Body.Close()

	board, err := h.BoardService.RenameBoard(r.Context(), boardID, req.NewName)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]*model.Board{
		"board": board,
	})
}

func (h *handler) getBoardWithColumnsAndCards(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")

	if boardID == "" {
		h.respondJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	board, err := h.BoardService.GetBoardWithColumnsAndCards(r.Context(), boardID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]*model.BoardWithColumns{
		"board": board,
	})
}

/* ====================== HELPER ====================== */

func (h *handler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
