package services

import (
	"anemone_backend-microservices/internal/kanban/model"
	"anemone_backend-microservices/internal/kanban/repository"
	"context"
)

type Service struct {
	BoardRepo  *repository.BoardRepo
	ColumnRepo *repository.ColumnRepo
	CardRepo   *repository.CardRepo
}

func NewService(br *repository.BoardRepo, cr *repository.ColumnRepo, cardRepo *repository.CardRepo) *Service {
	return &Service{BoardRepo: br, ColumnRepo: cr, CardRepo: cardRepo}
}

/* ====================== BOARD ====================== */

func (s *Service) CreateBoard(ctx context.Context, title string, userID int64) (*trello_model.Board, error) {
	return s.BoardRepo.CreateBoard(ctx, title, userID)
}

func (s *Service) GetOneUserBoard(ctx context.Context, boardID string) (*trello_model.BoardWithColumns, error) {
	return s.BoardRepo.GetOneUserBoard(ctx, boardID)
}

func (s *Service) GetAllUserBoards(ctx context.Context, userID int64) ([]*trello_model.Board, error) {
	return s.BoardRepo.GetAllUserBoards(ctx, userID)
}

func (s *Service) DeleteBoard(ctx context.Context, boardID string) error {
	return s.BoardRepo.DeleteBoard(ctx, boardID)
}

func (s *Service) RenameBoard(ctx context.Context, boardID string, newName string) (*trello_model.Board, error) {
	return s.BoardRepo.RenameBoard(ctx, boardID, newName)
}

func (s *Service) UpdateBoard(ctx context.Context, boardID string, userID int64, boardData []*trello_model.Column) error {
	return s.BoardRepo.UpdateBoard(ctx, boardID, userID, boardData)
}

/* ====================== COLUMN ====================== */

func (s *Service) CreateColumn(ctx context.Context, boardID, columnTitle string) (*trello_model.Column, error) {
	return s.ColumnRepo.CreateColumn(ctx, boardID, columnTitle)
}

func (s *Service) DeleteColumn(ctx context.Context, boardID, columnID string) error {
	return s.ColumnRepo.DeleteColumn(ctx, boardID, columnID)
}

func (s *Service) RenameColumn(ctx context.Context, boardID, columnID, newName string) (*trello_model.Column, error) {
	return s.ColumnRepo.RenameColumn(ctx, boardID, columnID, newName)
}

/* ====================== CARD ====================== */

func (s *Service) CreateCard(ctx context.Context, columnID, cardTitle string) (*trello_model.Card, error) {
	return s.CardRepo.CreateCard(ctx, columnID, cardTitle)
}

func (s *Service) DeleteCard(ctx context.Context, columnID, cardID string) error {
	return s.CardRepo.DeleteCard(ctx, columnID, cardID)
}

func (s *Service) RenameCard(ctx context.Context, columnID, cardID, newName string) (*trello_model.Card, error) {
	return s.CardRepo.RenameCard(ctx, columnID, cardID, newName)
}
