package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.GetContext(ctx,
		user,
		"SELECT id, email, hashed_password, salt, created_at, updated_at FROM users WHERE email = $1;", email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed get user by email [%s]: %w", email, err)
	}

	return user, nil
}

func (r *Repository) HasUser(ctx context.Context, email string) (bool, error) {
	if err := r.db.GetContext(ctx, &domain.User{}, "SELECT id FROM users WHERE email = $1", email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed check user on existence: %w", err)
	}

	return true, nil
}

// SaveUser insert given user in postgresql database.
func (r *Repository) SaveUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.NamedExecContext(ctx,
		"INSERT INTO users (hashed_password, email, salt) VALUES (:hashed_password, :email, :salt)", user)
	if err != nil {
		return fmt.Errorf("failed save user to postgresql database: %w", err)
	}

	return nil
}
