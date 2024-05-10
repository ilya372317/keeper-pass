package datarepo

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/jmoiron/sqlx"
)

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
