package datarepo

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (r *Repository) GetDataByID(ctx context.Context, id int) (domain.Data, error) {
	return domain.Data{}, nil
}
