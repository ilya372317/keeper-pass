package keyrepo

import (
	"context"
	"fmt"
)

func (r *Repository) SaveKey(ctx context.Context, key string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed begin transaction: %w", err)
	}

	if _, err = tx.ExecContext(ctx, "UPDATE keys SET is_current = false"); err != nil {
		return fmt.Errorf("failed make old keys not current: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO keys (key, is_current) VALUES ($1, true)", key)
	if err != nil {
		return fmt.Errorf("failed save key to storage: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed commit transaction: %w", err)
	}

	return nil
}
