package repository

import (
	"anemone_backend-microservices/internal/notes/model"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

/* ===================== NOTES ===================== */

func (r *Repo) CreateNote(ctx context.Context, n *model.Note) error {
	q := `
		INSERT INTO notes (user_id, title, content)
		VALUES ($1,$2,$3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, q, n.UserID, n.Title, n.Content,
	).Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)
}

func (r *Repo) GetNote(ctx context.Context, id int, userID int64) (*model.Note, error) {
	var n model.Note
	q := `
		SELECT id,user_id,title,content,is_deleted,folder_id,created_at,updated_at
		FROM notes
		WHERE id=$1 AND user_id=$2
	`
	err := r.db.GetContext(ctx, &n, q, id, userID)
	return &n, err
}

func (r *Repo) ListNotes(ctx context.Context, userID int64) ([]model.Note, error) {
	var notes []model.Note
	q := `
		SELECT id,title,content,is_deleted,folder_id,created_at,updated_at
		FROM notes
		WHERE user_id=$1
		ORDER BY updated_at DESC
	`
	err := r.db.SelectContext(ctx, &notes, q, userID)
	return notes, err
}

func (r *Repo) UpdateTitleByID(ctx context.Context, id int, newTitle string) (*model.Note, error) {
	q := `UPDATE notes SET title=$1, updated_at=NOW() WHERE id=$2 RETURNING *;`
	var updatedPage model.Note
	err := r.db.QueryRowContext(ctx, q, newTitle, id).Scan(&updatedPage.ID, &updatedPage.UserID, &updatedPage.Title, &updatedPage.Content, &updatedPage.IsDeleted, &updatedPage.FolderID, &updatedPage.UpdatedAt, &updatedPage.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, err

	}
	return &updatedPage, nil
}

func (r *Repo) UpdateNoteByID(
	ctx context.Context,
	noteID int,
	userID int64,
	newContent string,
) (*model.Note, error) {

	const q = `
		UPDATE notes
		SET content = $1,
		    updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING
			id, user_id, title, content, is_deleted,
			folder_id, created_at, updated_at
	`

	var n model.Note
	err := r.db.QueryRowContext(
		ctx, q, newContent, noteID, userID,
	).Scan(
		&n.ID,
		&n.UserID,
		&n.Title,
		&n.Content,
		&n.IsDeleted,
		&n.FolderID,
		&n.CreatedAt,
		&n.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("note not found")
	}
	if err != nil {
		return nil, err
	}

	return &n, nil
}

/* ===================== FOLDERS ===================== */

func (r *Repo) CreateFolder(ctx context.Context, f *model.Folder) error {
	const q = `
		INSERT INTO notes_folder (user_id, title)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		q,
		f.UserID,
		f.Title,
	).Scan(
		&f.ID,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
}

func (r *Repo) GetAllFolders(
	ctx context.Context,
	userID int64,
) ([]model.Folder, error) {

	const q = `
		SELECT
			id,
			user_id,
			title,
			created_at,
			updated_at
		FROM notes_folder
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var folders []model.Folder
	err := r.db.SelectContext(ctx, &folders, q, userID)
	if err != nil {
		return nil, err
	}

	return folders, nil
}

func (r *Repo) UpdateTitleFolder(
	ctx context.Context,
	folderID int,
	userID int64,
	newTitle string,
) (*model.Folder, error) {

	const q = `
		UPDATE notes_folder
		SET title = $1,
		    updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING
			id,
			user_id,
			title,
			created_at,
			updated_at
	`

	var f model.Folder
	err := r.db.QueryRowContext(
		ctx,
		q,
		newTitle,
		folderID,
		userID,
	).Scan(
		&f.ID,
		&f.UserID,
		&f.Title,
		&f.CreatedAt,
		&f.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("folder not found")
	}
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (r *Repo) DeleteFolderByID(
	ctx context.Context,
	folderID int,
	userID int64,
) error {

	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM notes_folder WHERE id = $1 AND user_id = $2`,
		folderID,
		userID,
	)

	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("folder not found")
	}

	return nil
}

/* ===================== NOTES <-> FOLDERS ===================== */

func (r *Repo) GetAllNotesFromFolder(
	ctx context.Context,
	folderID int,
	userID int64,
) ([]model.Note, error) {

	const q = `
		SELECT
			id, user_id, title, content, is_deleted,
			folder_id, created_at, updated_at
		FROM notes
		WHERE folder_id = $1
		  AND user_id = $2
		  AND is_deleted = false
		ORDER BY updated_at DESC
	`

	var notes []model.Note
	err := r.db.SelectContext(ctx, &notes, q, folderID, userID)
	return notes, err
}

func (r *Repo) AddNoteToFolder(
	ctx context.Context,
	noteID int,
	folderID int,
	userID int64,
) (*model.Note, error) {

	const q = `
		UPDATE notes
		SET folder_id = $1,
		    updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING
			id, user_id, title, content, is_deleted,
			folder_id, created_at, updated_at
	`

	var n model.Note
	err := r.db.QueryRowContext(
		ctx, q, folderID, noteID, userID,
	).Scan(
		&n.ID,
		&n.UserID,
		&n.Title,
		&n.Content,
		&n.IsDeleted,
		&n.FolderID,
		&n.CreatedAt,
		&n.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("note not found")
	}
	return &n, err
}

func (r *Repo) RemoveNoteFromFolder(
	ctx context.Context,
	noteID int,
	userID int64,
) (*model.Note, error) {

	const q = `
		UPDATE notes
		SET folder_id = NULL,
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $2
		RETURNING
			id, user_id, title, content, is_deleted,
			folder_id, created_at, updated_at
	`

	var n model.Note
	err := r.db.QueryRowContext(ctx, q, noteID, userID).
		Scan(
			&n.ID,
			&n.UserID,
			&n.Title,
			&n.Content,
			&n.IsDeleted,
			&n.FolderID,
			&n.CreatedAt,
			&n.UpdatedAt,
		)

	if err == sql.ErrNoRows {
		return nil, errors.New("note not found")
	}
	return &n, err
}

/* ===================== DELETED ===================== */

func (r *Repo) MarkDeletedNote(
	ctx context.Context,
	noteID int,
	userID int64,
) error {

	res, err := r.db.ExecContext(ctx,
		`UPDATE notes SET is_deleted=true, updated_at=NOW()
		  WHERE id=$1 AND user_id=$2`,
		noteID, userID,
	)

	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("note not found")
	}

	return nil
}

func (r *Repo) UnmarkDeletedNote(
	ctx context.Context,
	noteID int,
	userID int64,
) error {

	res, err := r.db.ExecContext(ctx,
		`UPDATE notes SET is_deleted=false, updated_at=NOW()
		  WHERE id=$1 AND user_id=$2`,
		noteID, userID,
	)

	if err != nil {
		return err
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("note not found")
	}

	return nil
}

func (r *Repo) DeleteAllMarkedNotes(ctx context.Context, userID int64) error {
	q := `DELETE FROM pages WHERE user_id=$1 AND is_deleted=true;`
	_, err := r.db.ExecContext(ctx, q, userID)

	if err != nil {
		return err
	}

	return nil
}

/* ===================== HELPERS ===================== */

func (r *Repo) IsNoteOwner(noteID int, userID int64) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1 FROM notes
			WHERE id = $1 AND user_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRow(q, noteID, userID).Scan(&exists)
	return exists, err
}

func (r *Repo) IsFolderOwner(folderID int, userID int64) (bool, error) {
	const q = `
		SELECT EXISTS(
			SELECT 1 FROM notes_folder
			WHERE id = $1 AND user_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRow(q, folderID, userID).Scan(&exists)
	return exists, err
}
