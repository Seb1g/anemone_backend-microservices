package model

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int            `db:"id" json:"id"`
	UserID    int64          `db:"user_id" json:"-"`
	Title     string         `db:"title" json:"title"`
	Content   string         `db:"content" json:"content"`
	IsDeleted bool           `db:"is_deleted" json:"is_deleted"`
	FolderID  sql.NullInt64  `db:"folder_id" json:"folder_id"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

type Folder struct {
	ID        int       `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"-"`
	Title     string    `db:"title" json:"title"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
