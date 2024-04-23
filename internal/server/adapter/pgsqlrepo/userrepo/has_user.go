package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (r *Repository) HasUser(ctx context.Context, email string) (bool, error) {
	if err := r.db.GetContext(ctx, &domain.User{}, "SELECT id FROM users WHERE email = $1", email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed check user on existence: %w", err)
	}

	return true, nil
}
