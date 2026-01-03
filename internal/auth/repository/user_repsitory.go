package repository

import (
	"anemone_backend-microservices/internal/auth/model"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(ctx context.Context, u *model.User) error {
	const q = `
		INSERT INTO users (email, password, is_verify)
		VALUES ($1, $2, false)
		RETURNING id, created_at
	`
	return r.DB.QueryRowContext(ctx, q, u.Email, u.Password).
		Scan(&u.ID, &u.CreatedAt)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const q = `
		SELECT id, email, password, created_at
		FROM users
		WHERE email = $1
	`

	var u model.User
	err := r.DB.QueryRowContext(ctx, q, email).
		Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	const q = `
		SELECT id, email, created_at
		FROM users
		WHERE id = $1
	`

	var u model.User
	err := r.DB.QueryRowContext(ctx, q, id).
		Scan(&u.ID, &u.Email, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) UpdatePassword(
	ctx context.Context,
	userID int64,
	newHash string,
) error {
	const q = `
		UPDATE users
		SET password = $1
		WHERE id = $2
	`
	_, err := r.DB.ExecContext(ctx, q, newHash, userID)
	return err
}
