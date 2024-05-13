package data

import (
	"context"
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

func (r *Repository) Delete(ctx context.Context, ids []int, userID uint, kinds []domain.Kind) error {
	query, args, err := sqlx.In("DELETE FROM data_records WHERE id IN (?) AND user_id = ? AND kind IN (?)",
		ids, userID, kinds)
	if err != nil {
		return fmt.Errorf("failed build query for delete data records: %w", err)
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed delete data records from postgresql database: %w", err)
	}

	return nil
}

func (r *Repository) GetDataByID(ctx context.Context, id int) (domain.Data, error) {
	var result domain.Data
	if err := r.db.GetContext(ctx, &result, "SELECT * FROM data_records WHERE id = $1", id); err != nil {
		return domain.Data{}, fmt.Errorf("failed get data record with id[%d]: %w", id, err)
	}

	return result, nil
}

// GetAll retrieve all data records from database.
func (r *Repository) GetAll(ctx context.Context, userID uint) ([]domain.Data, error) {
	result := make([]domain.Data, 0)
	err := r.db.SelectContext(ctx, &result, "SELECT * FROM data_records WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed get data_records from postgresql database: %w", err)
	}

	return result, nil
}

func (r *Repository) SaveData(ctx context.Context, data domain.Data) error {
	query := "INSERT INTO data_records " +
		"(payload, metadata, payload_nonce, crypto_key, crypto_key_nonce, kind, user_id)" +
		" VALUES (:payload, :metadata, :payload_nonce, :crypto_key, :crypto_key_nonce, :kind, :user_id)"
	if _, err := r.db.NamedExecContext(ctx, query, data); err != nil {
		return fmt.Errorf("failed save data to database: %w", err)
	}

	return nil
}

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
