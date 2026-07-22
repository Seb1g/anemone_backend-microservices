package service

import (
	"anemone_backend-kanban/internal/domain/model"
	"context"
	"fmt"
)

func (s *service) CreateBoard(ctx context.Context, title string, userID int64) (*model.Board, error) {
	if title == "" {
		return nil, fmt.Errorf("board title cannot be empty")
	}
	
	return s.BoardRepo.CreateBoard(ctx, title, userID)
}

func (s *service) GetAllUserBoards(ctx context.Context, userID int64) ([]*model.Board, error) {
	return s.BoardRepo.GetAllUserBoards(ctx, userID)
}

func (s *service) DeleteBoard(ctx context.Context, boardID string) error {
	return s.BoardRepo.DeleteBoard(ctx, boardID)
}

func (s *service) RenameBoard(ctx context.Context, boardID string, newName string) (*model.Board, error) {
	if newName == "" {
		return nil, fmt.Errorf("board title cannot be empty")
	}
	return s.BoardRepo.RenameBoard(ctx, boardID, newName)
}

func (s *service) GetBoardWithColumnsAndCards(ctx context.Context, boardID string) (*model.BoardWithColumns, error) {
	board, err := s.BoardRepo.GetBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}

	columns, err := s.GetAllColumnsFromBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}

	return &model.BoardWithColumns{
		ID:      board.ID,
		Title:   board.Title,
		Columns: columns,
	}, nil
}
