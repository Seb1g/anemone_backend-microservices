package service

import (
	"anemone_backend-kanban/internal/domain/model"
	"context"
)

type CardRepo interface {
	CreateCard(ctx context.Context, columnID, cardTitle string) (*model.Card, error)
	DeleteCard(ctx context.Context, columnID, cardID string) error
	RenameCard(ctx context.Context, columnID, cardID, newName string) (*model.Card, error)
	GetAllCardsFromColumn(ctx context.Context, columnID string) ([]*model.Card, error)
	UpdateCardPosition(ctx context.Context, cardID, fromColumnID, toColumnID string, newPosition float64) error
	
	GetCardPosition(ctx context.Context, cardID string) (float64, error)
}

type ColumnRepo interface {
	CreateColumn(ctx context.Context, boardID, columnTitle string) (*model.Column, error)
	DeleteColumn(ctx context.Context, boardID, columnID string) error
	RenameColumn(ctx context.Context, boardID, columnID, newName string) (*model.Column, error)
	GetAllColumnsFromBoard(ctx context.Context, boardID string) ([]*model.ColumnCardRow, error)
	UpdateColumnPosition(ctx context.Context, columnID, boardID string, newPosition float64) error

	GetColumnPosition(ctx context.Context, columnID string) (float64, error)
	GetBoardOwnerIDByColumnID(ctx context.Context, columnID string) (int64, error)
}

type BoardRepo interface {
	CreateBoard(ctx context.Context, title string, userID int64) (*model.Board, error)
	GetBoard(ctx context.Context, boardID string) (*model.Board, error)
	GetAllUserBoards(ctx context.Context, userID int64) ([]*model.Board, error)
	DeleteBoard(ctx context.Context, boardID string) error
	RenameBoard(ctx context.Context, boardID string, newName string) (*model.Board, error)
	
	GetBoardOwnerID(ctx context.Context, boardID string) (int64, error)
}