package key

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/jmoiron/sqlx"
)

// Repository for key storage.
type Repository struct {
	db *sqlx.DB
}

// New creates new repository.
func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetKey(ctx context.Context) (*domain.Keys, error) {
	var result domain.Keys

	err := r.db.GetContext(ctx, &result, "SELECT * FROM keys WHERE is_current = true")

	if err != nil {
		return nil, fmt.Errorf("failed get key: %w", err)
	}

	return &result, nil
}

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
