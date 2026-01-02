package services

import "context"

type TokenManager interface {
	NewPair(ctx context.Context, userID int64) (accessToken string, refreshToken string, err error)
	ParseAccess(token string) (int64, error)
	ParseRefresh(token string) (int64, error)
}
