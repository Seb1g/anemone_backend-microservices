package repo

import (
	"anemone_backend-kanban/internal/domain/apperr"
	"anemone_backend-kanban/internal/domain/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *repo) CreateBoard(ctx context.Context, title string, userID int64) (*model.Board, error) {
	board := &model.Board{}
	qBoard := `
		INSERT INTO boards (title, user_id)
		VALUES ($1, $2)
		RETURNING *;
	`

	err := r.DB.QueryRowxContext(ctx, qBoard, title, userID).StructScan(board)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperr.ErrBoardCreateFailed, err)
	}

	return board, nil
}

func (r *repo) GetBoard(ctx context.Context, boardID string) (*model.Board, error) {
	var board model.Board
	q := `SELECT * FROM boards WHERE id = $1;`

	err := r.DB.GetContext(ctx, &board, q, boardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrBoardNotFound
		}
		return nil, fmt.Errorf("%w: %v", apperr.ErrBoardGetFailed, err)
	}
	return &board, nil
}

func (r *repo) GetAllUserBoards(ctx context.Context, userID int64) ([]*model.Board, error) {
	var boards []*model.Board

	q := `SELECT * FROM boards WHERE user_id = $1 ORDER BY created_at DESC;`
	err := r.DB.SelectContext(ctx, &boards, q, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all user boards: %w", err)
	}

	return boards, nil
}

func (r *repo) DeleteBoard(ctx context.Context, boardID string) error {
	q := `DELETE FROM boards WHERE id = $1;`

	res, err := r.DB.ExecContext(ctx, q, boardID)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrBoardDeleteFailed, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return apperr.ErrBoardNotFound
	}
	return nil
}

func (r *repo) RenameBoard(ctx context.Context, boardID string, newName string) (*model.Board, error) {
	q := `UPDATE boards SET title = $1, updated_at = NOW() WHERE id = $2 RETURNING *;`
	var board model.Board

	err := r.DB.GetContext(ctx, &board, q, newName, boardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrBoardNotFound
		}
		return nil, fmt.Errorf("%w: %v", apperr.ErrBoardRenameFailed, err)
	}
	return &board, nil
}

// Helpers

func (r *repo) GetBoardOwnerID(ctx context.Context, boardID string) (int64, error) {
	var ownerID int64
	query := `SELECT user_id FROM boards WHERE id = $1;`

	err := r.DB.GetContext(ctx, &ownerID, query, boardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperr.ErrBoardNotFound
		}
		return 0, fmt.Errorf("%w: %v", apperr.ErrBoardGetFailed, err)
	}
	return ownerID, nil
}
