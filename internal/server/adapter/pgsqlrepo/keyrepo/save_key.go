package keyrepo

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (r *Repository) SaveKey(ctx context.Context, key *domain.Keys) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed begin transaction: %w", err)
	}

	if _, err = tx.ExecContext(ctx, "UPDATE keys SET is_current = false"); err != nil {
		return fmt.Errorf("failed make old keys not current: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO keys (key, is_current, nonce) VALUES ($1, true, $2)", key.Key, key.Nonce)
	if err != nil {
		return fmt.Errorf("failed save key to storage: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed commit transaction: %w", err)
	}

	return nil
}
