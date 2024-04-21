package userrepo

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user := domain.User{}
	err := r.db.GetContext(ctx,
		&user, "SELECT id, email, hashed_password, salt FROM users WHERE email = $1;", email)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed get user by email [%s]: %w", email, err)
	}

	return user, nil
}
