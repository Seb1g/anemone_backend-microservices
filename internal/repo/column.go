package repo

import (
	"anemone_backend-kanban/internal/domain/apperr"
	"anemone_backend-kanban/internal/domain/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *repo) CreateColumn(ctx context.Context, boardID, columnTitle string) (*model.Column, error) {
	q := `
		INSERT INTO columns (column_title, board_id, position)
		VALUES (
			$1,
			$2,
			(SELECT COALESCE(MAX(position), 0) + $3 FROM columns WHERE board_id = $2)
		) RETURNING *;
	`
	column := &model.Column{}

	err := r.DB.QueryRowxContext(ctx, q, columnTitle, boardID, r.DefaultStep).StructScan(column)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperr.ErrColumnCreateFailed, err)
	}

	return column, nil
}

func (r *repo) DeleteColumn(ctx context.Context, boardID, columnID string) error {
	q := `DELETE FROM columns WHERE id = $1 AND board_id = $2;`

	res, err := r.DB.ExecContext(ctx, q, columnID, boardID)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrColumnDeleteFailed, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return apperr.ErrColumnNotFound
	}

	return nil
}

func (r *repo) RenameColumn(ctx context.Context, boardID, columnID, newName string) (*model.Column, error) {
	q := `
		UPDATE columns
		SET column_title = $1
		WHERE id = $2 AND board_id = $3 RETURNING *;
	`
	var column model.Column

	err := r.DB.GetContext(ctx, &column, q, newName, columnID, boardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrColumnNotFound
		}
		return nil, fmt.Errorf("%w: %v", apperr.ErrColumnRenameFailed, err)
	}
	return &column, nil
}

func (r *repo) GetAllColumnsFromBoard(ctx context.Context, boardID string) ([]*model.ColumnCardRow, error) {
	q := `
		SELECT
			c.id AS column_id,
			c.column_title,
			c.board_id,
			c.position AS col_position,
			car.id AS card_id,
			car.content AS card_content,
			car.position AS card_position
		FROM columns c
		LEFT JOIN cards car ON c.id = car.column_id
		WHERE c.board_id = $1
		ORDER BY c.position ASC, car.position ASC;
	`
	var rows []*model.ColumnCardRow
	err := r.DB.SelectContext(ctx, &rows, q, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns with cards: %w", err)
	}

	return rows, nil
}

func (r *repo) UpdateColumnPosition(ctx context.Context, columnID, boardID string, newPosition float64) error {
	q := `
		UPDATE columns
		SET position = $1
		WHERE id = $2 AND board_id = $3;
	`
	res, err := r.DB.ExecContext(ctx, q, newPosition, columnID, boardID)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrColumnMoveFailed, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return apperr.ErrColumnNotFound
	}

	return nil
}

// Helpers

func (r *repo) GetColumnPosition(ctx context.Context, columnID string) (float64, error) {
	var pos float64
	q := `SELECT position FROM columns WHERE id = $1;`

	err := r.DB.GetContext(ctx, &pos, q, columnID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperr.ErrColumnNotFound
		}
		return 0, fmt.Errorf("%w: %v", apperr.ErrInternal, err)
	}
	return pos, nil
}

func (r *repo) GetBoardOwnerIDByColumnID(ctx context.Context, columnID string) (int64, error) {
	var ownerID int64
	query := `
		SELECT b.user_id
		FROM boards b
		JOIN columns c ON b.id = c.board_id
		WHERE c.id = $1;
	`

	err := r.DB.GetContext(ctx, &ownerID, query, columnID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperr.ErrBoardNotFound
		}
		return 0, fmt.Errorf("%w: %v", apperr.ErrBoardGetFailed, err)
	}
	return ownerID, nil
}
