package keyrepo

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (r *Repository) GetKey(ctx context.Context) (*domain.Keys, error) {
	var result domain.Keys

	err := r.db.GetContext(ctx, &result, "SELECT * FROM keys WHERE is_current = true")

	if err != nil {
		return nil, fmt.Errorf("failed get key: %w", err)
	}

	return &result, nil
}
