package datarepo

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (r *Repository) GetDataByID(ctx context.Context, id int) (domain.Data, error) {
	var result domain.Data
	if err := r.db.GetContext(ctx, &result, "SELECT * FROM data_records WHERE id = $1", id); err != nil {
		return domain.Data{}, fmt.Errorf("failed get data record with id[%d]: %w", id, err)
	}

	return result, nil
}
