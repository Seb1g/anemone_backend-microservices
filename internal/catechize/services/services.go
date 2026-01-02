package services

import (
	"anemone_backend-microservices/internal/catechize/model"
	"anemone_backend-microservices/internal/catechize/repository"
	"context"
)

type Service struct {
	repo *repository.Repo
}

func NewService(r *repository.Repo) *Service {
	return &Service{repo: r}
}

func (s *Service) AddResult(ctx context.Context, userID int64, req model.QuizResult) (*model.QuizResult, error) {
	req.UserID = userID
	if err := s.repo.AddResult(ctx, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *Service) GetResults(ctx context.Context, userID int64) ([]model.QuizResult, error) {
	return s.repo.ListResults(ctx, userID)
}

func (s *Service) ClearResultsUser(ctx context.Context, userID int64) error {
	return s.repo.ClearResultsUser(ctx, userID)
}
