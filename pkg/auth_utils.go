// Package pkg contains re-use utils
package pkg

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(tokenString string, secret string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id missing")
	}

	return int64(id), nil
}
