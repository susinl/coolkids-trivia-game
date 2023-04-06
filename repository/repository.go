package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewConnection(connectionString string) *sql.DB {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CheckGameCode(gameCode string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT game_code FROM participant WHERE game_code = $1)", gameCode).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
