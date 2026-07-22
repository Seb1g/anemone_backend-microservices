package service

import (
	"anemone_backend-kanban/internal/domain/apperr"
	"anemone_backend-kanban/internal/domain/model"
	"context"
	"errors"
	"fmt"
)

func (s *service) CreateCard(ctx context.Context, columnID, cardTitle string) (*model.CardResponse, error) {
	card, err := s.CardRepo.CreateCard(ctx, columnID, cardTitle)
	if err != nil {
		return nil, err
	}

	return &model.CardResponse{
		ID:       card.ID,
		Content:  card.Content,
		Position: card.Position,
		ColumnID: card.ColumnID,
	}, nil
}

func (s *service) DeleteCard(ctx context.Context, columnID, cardID string) error {
	err := s.CardRepo.DeleteCard(ctx, columnID, cardID)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) RenameCard(ctx context.Context, columnID, cardID, newName string) (*model.Card, error) {
	if newName == "" {
		return nil, fmt.Errorf("column title cannot be empty")
	}

	card, err := s.CardRepo.RenameCard(ctx, columnID, cardID, newName)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (s *service) GetAllCardsFromColumn(ctx context.Context, columnID string) ([]*model.Card, error) {
	cards, err := s.CardRepo.GetAllCardsFromColumn(ctx, columnID)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (s *service) MoveCard(ctx context.Context, params model.MoveCardParams) error {
	var newPosition float64

	if params.PrevCardID == nil && params.NextCardID == nil {
		newPosition = s.DefaultStep
	} else {
		var prevPos, nextPos *float64

		if params.PrevCardID != nil {
			p, err := s.CardRepo.GetCardPosition(ctx, *params.PrevCardID)
			if err != nil {
				if errors.Is(err, apperr.ErrCardNotFound) {
					return fmt.Errorf("invalid neighbor: previous card not found")
				}
				return fmt.Errorf("failed to get prev card position: %w", err)
			}
			prevPos = &p
		}

		if params.NextCardID != nil {
			n, err := s.CardRepo.GetCardPosition(ctx, *params.NextCardID)
			if err != nil {
				if errors.Is(err, apperr.ErrCardNotFound) {
					return fmt.Errorf("invalid neighbor: next card not found")
				}
				return fmt.Errorf("failed to get next card position: %w", err)
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

	return s.CardRepo.UpdateCardPosition(
		ctx,
		params.CardID,
		params.FromColumnID,
		params.ToColumnID,
		newPosition,
	)
}
