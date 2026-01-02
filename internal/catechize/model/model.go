package model

import "time"

type QuizResult struct {
	ID                  int64     `db:"id" json:"id"`
	UserID              int64     `db:"user_id" json:"-"`
	CurrentAnswers      string    `db:"current_answers" json:"current_answer"`
	CountQuestions      string    `db:"count_questions" json:"count_questions"`
	TypeQuestions       string    `db:"type_questions" json:"type_questions"`
	DifficultyQuestions string    `db:"difficulty_questions" json:"difficulty_questions"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
}
