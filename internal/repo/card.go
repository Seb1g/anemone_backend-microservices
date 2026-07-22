package repo

import (
	"anemone_backend-kanban/internal/domain/apperr"
	"anemone_backend-kanban/internal/domain/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *repo) CreateCard(ctx context.Context, columnID, cardTitle string) (*model.Card, error) {
	card := &model.Card{}

	q := `
		INSERT INTO cards (content, column_id, position)
		VALUES (
			$1,
			$2,
			(SELECT COALESCE(MAX(position), 0) + $3 FROM cards WHERE column_id = $2)
		)
		RETURNING *;
	`

	err := r.DB.QueryRowxContext(ctx, q, cardTitle, columnID, r.DefaultStep).StructScan(card)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperr.ErrCardCreateFailed, err)
	}

	return card, nil
}

func (r *repo) DeleteCard(ctx context.Context, columnID, cardID string) error {
	q := `DELETE FROM cards WHERE id = $1 AND column_id = $2;`

	res, err := r.DB.ExecContext(ctx, q, cardID, columnID)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrCardDeleteFailed, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return apperr.ErrCardNotFound
	}

	return nil
}

func (r *repo) RenameCard(ctx context.Context, columnID, cardID, newName string) (*model.Card, error) {
	q := `
		UPDATE cards
		SET content = $1
		WHERE id = $2 AND column_id = $3 RETURNING *;
	`
	var card model.Card

	err := r.DB.GetContext(ctx, &card, q, newName, cardID, columnID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrCardNotFound
		}
		return nil, fmt.Errorf("%w: %v", apperr.ErrCardRenameFailed, err)
	}

	return &card, nil
}

func (r *repo) GetAllCardsFromColumn(ctx context.Context, columnID string) ([]*model.Card, error) {
	q := `SELECT * FROM cards WHERE column_id = $1 ORDER BY position ASC;`
	var cards []*model.Card
	
	err := r.DB.SelectContext(ctx, &cards, q, columnID)
	if err != nil {
		return nil, fmt.Errorf("failed to select cards for column: %w", err)
	}

	return cards, nil
}

func (r *repo) UpdateCardPosition(ctx context.Context, cardID, fromColumnID, toColumnID string, newPosition float64) error {
	q := `
		UPDATE cards
		SET column_id = $1, position = $2
		WHERE id = $3 AND column_id = $4;
	`
	res, err := r.DB.ExecContext(ctx, q, toColumnID, newPosition, cardID, fromColumnID)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrCardMoveFailed, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return apperr.ErrCardNotFound
	}

	return nil
}

// Helpers

func (r *repo) GetCardPosition(ctx context.Context, cardID string) (float64, error) {
	var pos float64
	q := `SELECT position FROM cards WHERE id = $1;`

	err := r.DB.GetContext(ctx, &pos, q, cardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperr.ErrCardNotFound
		}
		return 0, fmt.Errorf("%w: %v", apperr.ErrInternal, err)
	}
	return pos, nil
}
