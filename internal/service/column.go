package service

import (
	"anemone_backend-kanban/internal/domain/apperr"
	"anemone_backend-kanban/internal/domain/model"
	"context"
	"errors"
	"fmt"
)

func (s *service) CreateColumn(ctx context.Context, boardID, columnTitle string) (*model.Column, error) {
	column, err := s.ColumnRepo.CreateColumn(ctx, boardID, columnTitle)
	if err != nil {
		return nil, err
	}

	return &model.Column{
		ID:       column.ID,
		Title:    column.Title,
		BoardID:  column.BoardID,
		Position: column.Position,
		Cards:    column.Cards,
	}, nil
}

func (s *service) DeleteColumn(ctx context.Context, boardID, columnID string) error {
	err := s.ColumnRepo.DeleteColumn(ctx, boardID, columnID)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) RenameColumn(ctx context.Context, boardID, columnID, newName string) (*model.Column, error) {
	if newName == "" {
		return nil, fmt.Errorf("column title cannot be empty")
	}

	column, err := s.ColumnRepo.RenameColumn(ctx, boardID, columnID, newName)
	if err != nil {
		return nil, err
	}

	return column, nil
}

func (s *service) GetAllColumnsFromBoard(ctx context.Context, boardID string) ([]*model.Column, error) {
	columnsRows, err := s.ColumnRepo.GetAllColumnsFromBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}

	columnsMap := make(map[string]*model.Column)
	var columns []*model.Column

	for _, row := range columnsRows {
		col, exists := columnsMap[row.ColumnID]
		if !exists {
			col = &model.Column{
				ID:       row.ColumnID,
				Title:    row.ColumnTitle,
				BoardID:  row.BoardID,
				Position: row.ColPosition,
				Cards:    make([]*model.Card, 0),
			}
			columnsMap[row.ColumnID] = col
			columns = append(columns, col)
		}

		if row.CardID != nil {
			card := &model.Card{
				ID:       *row.CardID,
				Content:  *row.CardContent,
				ColumnID: row.ColumnID,
				Position: *row.CardPosition,
			}
			col.Cards = append(col.Cards, card)
		}
	}

	return columns, nil
}

func (s *service) MoveColumn(ctx context.Context, params model.MoveColumnParams) error {
	var newPosition float64

	if params.PrevColumnID == nil && params.NextColumnID == nil {
		newPosition = s.DefaultStep
	} else {
		var prevPos, nextPos *float64

		if params.PrevColumnID != nil {
			p, err := s.ColumnRepo.GetColumnPosition(ctx, *params.PrevColumnID)
			if err != nil {
				if errors.Is(err, apperr.ErrColumnNotFound) {
					return fmt.Errorf("invalid neighbor: previous card not found")
				}
				return fmt.Errorf("failed to get prev column position: %w", err)
			}
			prevPos = &p
		}

		if params.NextColumnID != nil {
			n, err := s.ColumnRepo.GetColumnPosition(ctx, *params.NextColumnID)
			if err != nil {
				if errors.Is(err, apperr.ErrColumnNotFound) {
					return fmt.Errorf("invalid neighbor: next column not found")
				}
				return fmt.Errorf("failed to get next column position: %w", err)
			}
			nextPos = &n
		}

		if prevPos != nil && nextPos != nil {
			newPosition = (*prevPos + *nextPos) / 2.0
		} else if prevPos != nil {
			newPosition = *prevPos + s.DefaultStep
		} else if nextPos != nil {
			newPosition = *nextPos / 2.0
		}
	}

	return s.ColumnRepo.UpdateColumnPosition(
		ctx,
		params.ColumnID,
		params.BoardID,
		newPosition,
	)
}
