package userrepo

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

// SaveUser insert given user in postgresql database.
func (r *Repository) SaveUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.NamedExecContext(ctx,
		"INSERT INTO users (hashed_password, email, salt) VALUES (:hashed_password, :email, :salt)", user)
	if err != nil {
		return fmt.Errorf("failed save user to postgresql database: %w", err)
	}

	return nil
}
