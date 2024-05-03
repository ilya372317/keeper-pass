package data

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) GetAndDecryptData(ctx context.Context, id int64) (*domain.Data, error) {
	return &domain.Data{}, nil
}
