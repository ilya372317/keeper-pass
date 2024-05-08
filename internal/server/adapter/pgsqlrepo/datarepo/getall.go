package datarepo

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

// GetAll retrieve all data records from database.
func (r *Repository) GetAll(ctx context.Context) ([]domain.Data, error) {
	result := make([]domain.Data, 0)
	if err := r.db.SelectContext(ctx, &result, "SELECT * FROM data_records"); err != nil {
		return nil, fmt.Errorf("failed get data_records from postgresql database: %w", err)
	}

	return result, nil
}
