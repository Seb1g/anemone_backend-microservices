package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenManager struct {
	accessSecret  []byte
	refreshSecret []byte

	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTTokenManager(
	accessSecret string,
	refreshSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *JWTTokenManager {
	return &JWTTokenManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (m *JWTTokenManager) NewPair(
	_ context.Context,
	userID int64,
) (string, string, error) {

	now := time.Now()

	accessClaims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString(m.accessSecret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString(m.refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *JWTTokenManager) ParseAccess(tokenStr string) (int64, error) {
	return m.parse(tokenStr, m.accessSecret)
}

func (m *JWTTokenManager) ParseRefresh(tokenStr string) (int64, error) {
	return m.parse(tokenStr, m.refreshSecret)
}

func (m *JWTTokenManager) parse(tokenStr string, secret []byte) (int64, error) {
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
