package services

import (
	"context"

	"anemone_backend-microservices/internal/notes/model"
	"anemone_backend-microservices/internal/notes/repository"
)

type Service struct {
	repo *repository.Repo
}

func NewService(repo *repository.Repo) *Service {
	return &Service{repo: repo}
}

/* ===================== NOTES ===================== */

func (s *Service) CreateNote(
	ctx context.Context,
	userID int64,
	title string,
	content string,
) (*model.Note, error) {

	note := &model.Note{
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	if err := s.repo.CreateNote(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *Service) GetNote(
	ctx context.Context,
	noteID int,
	userID int64,
) (*model.Note, error) {
	return s.repo.GetNote(ctx, noteID, userID)
}

func (s *Service) ListNotes(
	ctx context.Context,
	userID int64,
) ([]model.Note, error) {
	return s.repo.ListNotes(ctx, userID)
}

func (s *Service) UpdateTitle(
	ctx context.Context,
	noteID int,
	userID int64,
	newTitle string,
) (*model.Note, error) {
	return s.repo.UpdateTitleByID(ctx, noteID, newTitle)
}

func (s *Service) UpdateContent(
	ctx context.Context,
	noteID int,
	userID int64,
	newContent string,
) (*model.Note, error) {
	return s.repo.UpdateNoteByID(ctx, noteID, userID, newContent)
}

/* ===================== FOLDERS ===================== */

func (s *Service) CreateFolder(
	ctx context.Context,
	userID int64,
	title string,
) (*model.Folder, error) {

	folder := &model.Folder{
		UserID: userID,
		Title:  title,
	}

	if err := s.repo.CreateFolder(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *Service) GetAllFolders(
	ctx context.Context,
	userID int64,
) ([]model.Folder, error) {
	return s.repo.GetAllFolders(ctx, userID)
}

func (s *Service) UpdateTitleFolder(
	ctx context.Context,
	folderID int,
	userID int64,
	newTitle string,
) (*model.Folder, error) {
	return s.repo.UpdateTitleFolder(ctx, folderID, userID, newTitle)
}

func (s *Service) DeleteFolder(
	ctx context.Context,
	folderID int,
	userID int64,
) error {
	return s.repo.DeleteFolderByID(ctx, folderID, userID)
}

/* ===================== NOTES <-> FOLDERS ===================== */

func (s *Service) GetAllNotesFromFolder(
	ctx context.Context,
	folderID int,
	userID int64,
) ([]model.Note, error) {
	return s.repo.GetAllNotesFromFolder(ctx, folderID, userID)
}

func (s *Service) AddNoteToFolder(
	ctx context.Context,
	noteID int,
	folderID int,
	userID int64,
) (*model.Note, error) {
	return s.repo.AddNoteToFolder(ctx, noteID, folderID, userID)
}

func (s *Service) RemoveNoteFromFolder(
	ctx context.Context,
	noteID int,
	userID int64,
) (*model.Note, error) {
	return s.repo.RemoveNoteFromFolder(ctx, noteID, userID)
}

/* ===================== DELETED ===================== */

func (s *Service) MarkDeletedNote(
	ctx context.Context,
	noteID int,
	userID int64,
) error {
	return s.repo.MarkDeletedNote(ctx, noteID, userID)
}

func (s *Service) UnmarkDeletedNote(
	ctx context.Context,
	noteID int,
	userID int64,
) error {
	return s.repo.UnmarkDeletedNote(ctx, noteID, userID)
}

func (s *Service) DeleteAllMarkedNotes(
	ctx context.Context,
	userID int64,
) error {
	return s.repo.DeleteAllMarkedNotes(ctx, userID)
}
