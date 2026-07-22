// Package repo contains istanse to db
package repo

import "github.com/jmoiron/sqlx"

type repo struct {
	DB          *sqlx.DB
	DefaultStep float64
}

func NewRepo(db *sqlx.DB, df float64) *repo {
	return &repo{
		DB: db,
		DefaultStep: df,
	}
}
