package datarepo

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (r *Repository) SaveData(ctx context.Context, data domain.Data) error {
	return nil
}

func (r *Repository) GetData(ctx context.Context, id int) (*domain.Data, error) {
	return &domain.Data{}, nil
}
