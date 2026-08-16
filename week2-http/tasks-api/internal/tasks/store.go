package tasks

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(databaseURL string) (*TaskStore, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &TaskStore{db: db}, nil
}
