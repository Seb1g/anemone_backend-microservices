package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type RefreshRepo struct {
	DB *sqlx.DB
}

func NewRefreshRepo(db *sqlx.DB) *RefreshRepo {
	return &RefreshRepo{DB: db}
}

func (r *RefreshRepo) Store(
	ctx context.Context,
	userID int64,
	token string,
	exp time.Time,
) error {
	const q = `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.DB.ExecContext(ctx, q, userID, token, exp)
	return err
}

func (r *RefreshRepo) Exists(
	ctx context.Context,
	userID int64,
	token string,
) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1
			FROM refresh_tokens
			WHERE user_id = $1
			AND token = $2
			AND expires_at > NOW()
		)
	`

	var exists bool
	err := r.DB.QueryRowContext(ctx, q, userID, token).Scan(&exists)
	return exists, err
}

func (r *RefreshRepo) Delete(
	ctx context.Context,
	userID int64,
	token string,
) error {
	const q = `
		DELETE FROM refresh_tokens
		WHERE user_id = $1 AND token = $2
	`
	_, err := r.DB.ExecContext(ctx, q, userID, token)
	return err
}
