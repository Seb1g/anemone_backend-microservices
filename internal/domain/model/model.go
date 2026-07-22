// Package model
package model

import (
	"time"
)

// Card Repo&Req/Res
type Card struct {
	ID       string  `db:"id"`
	Content  string  `db:"content"`
	ColumnID string  `db:"column_id"`
	Position float64 `db:"position"`
}

type MoveCardParams struct {
	CardID       string
	FromColumnID string
	ToColumnID   string
	PrevCardID   *string
	NextCardID   *string
}

type CardResponse struct {
	ID       string  `json:"id"`
	Content  string  `json:"content"`
	ColumnID string  `json:"column_id"`
	Position float64 `json:"position"`
}

// Column Repo&Req/Res
type Column struct {
	ID       string  `db:"id" json:"id"`
	Title    string  `db:"column_title" json:"title"`
	BoardID  string  `db:"board_id" json:"-"`
	Position float64 `db:"position" json:"position"`
	Cards    []*Card `db:"cards" json:"cards"`
}

type MoveColumnParams struct {
	ColumnID     string
	BoardID      string
	PrevColumnID *string
	NextColumnID *string
}

type ColumnCardRow struct {
	ColumnID     string   `db:"id"`
	ColumnTitle  string   `db:"column_title"`
	BoardID      string   `db:"board_id"`
	ColPosition  float64  `db:"position"`
	CardID       *string  `db:"id"`
	CardContent  *string  `db:"content"`
	CardPosition *float64 `db:"position"`
}

type Board struct {
	ID        string     `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	UserID    int64      `db:"user_id" json:"user_id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type BoardWithColumns struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Columns []*Column `json:"columns"`
}
