package userrepo

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return domain.User{}, nil
}
