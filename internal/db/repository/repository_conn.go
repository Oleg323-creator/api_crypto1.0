package repository

import "database/sql"

type Repository struct {
	DB *sql.DB
}

func NewRepository(dbConn *sql.DB) *Repository {
	return &Repository{DB: dbConn}
}
