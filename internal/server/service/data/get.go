package data

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) GetAndDecryptData(ctx context.Context, id int64) (*domain.Data, error) {
	return &domain.Data{}, nil
}
func (s *Service) EncryptAndUpdateData(ctx context.Context, d dto.UpdateSimpleDataDTO) error {
	return nil
}
