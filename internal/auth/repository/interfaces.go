package repository

import (
	"anemone_backend-microservices/internal/auth/model"
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	UpdatePassword(ctx context.Context, userID int64, newHash string) error
}

type RefreshRepository interface {
	Store(ctx context.Context, userID int64, token string, exp time.Time) error
	Exists(ctx context.Context, userID int64, token string) (bool, error)
	Delete(ctx context.Context, userID int64, token string) error
}
