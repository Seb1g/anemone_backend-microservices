package services

import (
	"anemone_backend-microservices/internal/auth/model"
	"anemone_backend-microservices/internal/auth/repository"
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repository.UserRepository
	refreshRepo repository.RefreshRepository
	tokens      TokenManager
}

func NewAuthService(
	u repository.UserRepository,
	r repository.RefreshRepository,
	t TokenManager,
) *AuthService {
	return &AuthService{
		userRepo:    u,
		refreshRepo: r,
		tokens:      t,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (string, string, *model.User, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", nil, errors.New("failed to hash password")
	}

	u := &model.User{Email: email, Password: string(hash)}
	if err := s.userRepo.Create(ctx, u); err != nil {
		return "", "", nil, err
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, u)

	if err != nil {
		return "", "", nil, err
	}
	return accessToken, refreshToken, u, nil
}

func (s *AuthService) generateTokens(
	ctx context.Context,
	user *model.User,
) (string, string, error) {

	access, refresh, err := s.tokens.NewPair(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	exp := time.Now().Add(7 * 24 * time.Hour)
	if err := s.refreshRepo.Store(ctx, user.ID, refresh, exp); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (string, string, *model.User, error) {

	userID, err := s.tokens.ParseRefresh(refreshToken)
	if err != nil {
		return "", "", nil, errors.New("invalid refresh token")
	}

	ok, err := s.refreshRepo.Exists(ctx, userID, refreshToken)
	if err != nil || !ok {
		return "", "", nil, errors.New("refresh token not found")
	}

	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", nil, errors.New("user not found")
	}

	access, refresh, err := s.generateTokens(ctx, u)
	if err != nil {
		return "", "", nil, err
	}

	_ = s.refreshRepo.Delete(ctx, userID, refreshToken)

	return access, refresh, u, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, *model.User, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return "", "", nil, errors.New("invalid credentials")
	}
	access, refresh, err := s.generateTokens(ctx, u)
	if err != nil {
		return "", "", nil, err
	}
	return access, refresh, u, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	userID, err := s.tokens.ParseRefresh(refreshToken)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	return s.refreshRepo.Delete(ctx, userID, refreshToken)
}

func (s *AuthService) ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error {
	email = strings.TrimSpace(email)
	oldPassword = strings.TrimSpace(oldPassword)
	newPassword = strings.TrimSpace(newPassword)

	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword)) != nil {
		return errors.New("invalid old password")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	return s.userRepo.UpdatePassword(ctx, u.ID, string(newHash))
}

func (s *AuthService) ResetPassword(ctx context.Context, email, newPassword string) error {
	email = strings.TrimSpace(email)
	newPassword = strings.TrimSpace(newPassword)

	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("internal server error")
	}
	return s.userRepo.UpdatePassword(ctx, u.ID, string(hash))
}

func (s *AuthService) ParseAccessToken(tokenStr string) (int64, error) {
	return s.tokens.ParseAccess(tokenStr)
}
