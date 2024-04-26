package keyrepo

import "github.com/jmoiron/sqlx"

// Repository for key storage.
type Repository struct {
	db *sqlx.DB
}

// New creates new repository.
func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}
