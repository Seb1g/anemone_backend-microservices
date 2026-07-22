package api

import (
	"anemone_backend-kanban/internal/domain/model"
	"context"
)

type CardServiceInterface interface {
	CreateCard(ctx context.Context, columnID, cardTitle string) (*model.CardResponse, error)
	DeleteCard(ctx context.Context, columnID, cardID string) error
	RenameCard(ctx context.Context, columnID, cardID, newName string) (*model.Card, error)
	GetAllCardsFromColumn(ctx context.Context, columnID string) ([]*model.Card, error)
	MoveCard(ctx context.Context, params model.MoveCardParams) error
}

type ColumnServiceInterface interface {
	CreateColumn(ctx context.Context, boardID, columnTitle string) (*model.Column, error)
	DeleteColumn(ctx context.Context, boardID, columnID string) error
	RenameColumn(ctx context.Context, boardID, columnID, newName string) (*model.Column, error)
	GetAllColumnsFromBoard(ctx context.Context, boardID string) ([]*model.Column, error)
	MoveColumn(ctx context.Context, params model.MoveColumnParams) error
}

type BoardServiceInterface interface {
	CreateBoard(ctx context.Context, title string, userID int64) (*model.Board, error)
	GetAllUserBoards(ctx context.Context, userID int64) ([]*model.Board, error)
	DeleteBoard(ctx context.Context, boardID string) error
	RenameBoard(ctx context.Context, boardID string, newName string) (*model.Board, error)
	GetBoardWithColumnsAndCards(ctx context.Context, boardID string) (*model.BoardWithColumns, error)
}
