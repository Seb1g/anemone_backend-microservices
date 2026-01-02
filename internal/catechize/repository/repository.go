package repository

import (
	"anemone_backend-microservices/internal/catechize/model"
	"context"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) AddResult(ctx context.Context, u *model.QuizResult) error {
	const q = `
		INSERT INTO quiz (
			user_id,
			current_answers,
			count_questions,
			type_questions,
			difficulty_questions
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		q,
		u.UserID,
		u.CurrentAnswers,
		u.CountQuestions,
		u.TypeQuestions,
		u.DifficultyQuestions,
	).Scan(&u.ID, &u.CreatedAt)
}

func (r *Repo) ListResults(ctx context.Context, userID int64) ([]model.QuizResult, error) {
	var results []model.QuizResult

	const q = `
		SELECT
			id,
			user_id,
			current_answers,
			count_questions,
			type_questions,
			difficulty_questions,
			created_at
		FROM quiz
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &results, q, userID)
	return results, err
}

func (r *Repo) ClearResultsUser(ctx context.Context, userID int64) error {
	const q = `
		DELETE FROM quiz WHERE user_id = $1
	`

	_, err := r.db.ExecContext(ctx, q, userID)
	return err
}
