package datarepo

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (r *Repository) UpdateByID(ctx context.Context, id int, dto domain.Data) error {
	if _, err := r.db.ExecContext(
		ctx,
		"UPDATE data_records SET payload = $1, metadata = $2, payload_nonce = $3,"+
			" crypto_key = $4, crypto_key_nonce = $5 WHERE id = $6",
		dto.Payload, dto.Metadata, dto.PayloadNonce, dto.CryptoKey, dto.CryptoKeyNonce, id,
	); err != nil {
		return fmt.Errorf("failed update data_records: %w", err)
	}

	return nil
}
