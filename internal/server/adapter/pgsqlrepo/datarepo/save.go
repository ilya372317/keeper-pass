package datarepo

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (r *Repository) SaveData(ctx context.Context, data domain.Data) error {
	query := "INSERT INTO data_records " +
		"(payload, metadata, payload_nonce, crypto_key, crypto_key_nonce, kind, user_id)" +
		" VALUES (:payload, :metadata, :payload_nonce, :crypto_key, :crypto_key_nonce, :kind, :user_id)"
	if _, err := r.db.NamedExecContext(ctx, query, data); err != nil {
		return fmt.Errorf("failed save data to database: %w", err)
	}

	return nil
}
