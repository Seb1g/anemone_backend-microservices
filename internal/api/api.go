// Package api
package api

import (
	"net/http"
)

type handler struct {
	CardService   CardServiceInterface
	ColumnService ColumnServiceInterface
	BoardService  BoardServiceInterface
	AccessSecret  string
}

func NewHandler(
	cardService CardServiceInterface,
	columnService ColumnServiceInterface,
	boardService BoardServiceInterface,
	accessSecret string,
) *handler {
	return &handler{
		CardService:   cardService,
		ColumnService: columnService,
		BoardService:  boardService,
		AccessSecret:  accessSecret,
	}
}

type Middleware func(http.Handler) http.Handler

func Adapt(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (h *handler) Routes(mux *http.ServeMux) {
	am := AuthMiddleware(h.AccessSecret)
	
	// --- CARDS ---
	mux.Handle(
		"POST /api/v1/kanban/columns/{column_id}/cards", 
		Adapt(http.HandlerFunc(h.createCard), am),
	)
	mux.Handle(
		"DELETE /api/v1/kanban/columns/{column_id}/cards/{card_id}",
		Adapt(http.HandlerFunc(h.deleteCard), am),
	)
	mux.Handle(
		"PUT /api/v1/kanban/columns/{column_id}/cards/{card_id}",
		Adapt(http.HandlerFunc(h.renameCard), am),
	)
	mux.Handle(
		"GET /api/v1/kanban/columns/{column_id}/cards",
		Adapt(http.HandlerFunc(h.getAllCardsFromColumn), am),
	)
	mux.Handle(
		"PUT /api/v1/kanban/move/columns/{column_id}/cards/{card_id}",
		Adapt(http.HandlerFunc(h.moveCard), am),
	)

	// --- COLUMNS ---
	mux.Handle(
		"POST /api/v1/kanban/boards/{board_id}/columns",
		Adapt(http.HandlerFunc(h.createColumn), am),
	)
	mux.Handle(
		"DELETE /api/v1/kanban/boards/{board_id}/columns/{column_id}",
		Adapt(http.HandlerFunc(h.deleteColumn), am),
	)
	mux.Handle(
		"PUT /api/v1/kanban/boards/{board_id}/columns/{column_id}",
		Adapt(http.HandlerFunc(h.renameColumn), am),
	)
	mux.Handle(
		"GET /api/v1/kanban/boards/{board_id}/columns",
		Adapt(http.HandlerFunc(h.getAllColumnsFromBoard), am),
	)
	mux.Handle(
		"PUT /api/v1/kanban/move/board/{board_id}/column/{column_id}",
		Adapt(http.HandlerFunc(h.moveColumn), am),
	)

	// --- BOARDS ---
	mux.Handle(
		"POST /api/v1/kanban/boards",
		Adapt(http.HandlerFunc(h.createBoard), am),
	)
	mux.Handle(
		"GET /api/v1/kanban/boards",
		Adapt(http.HandlerFunc(h.getAllUserBoards), am),
	)
	mux.Handle(
		"GET /api/v1/kanban/boards/{board_id}",
		Adapt(http.HandlerFunc(h.getBoardWithColumnsAndCards), am),
	)
	mux.Handle(
		"DELETE /api/v1/kanban/boards/{board_id}",
		Adapt(http.HandlerFunc(h.deleteBoard), am),
	)
	mux.Handle(
		"PUT /api/v1/kanban/boards/{board_id}",
		Adapt(http.HandlerFunc(h.renameBoard), am),
	)
}