package data

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) GetAndDecryptData(ctx context.Context, id int64) (domain.Data, error) {
	data, err := s.dataStorage.GetDataByID(ctx, int(id))
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed get data from storage: %w", err)
	}
	decryptedData, err := s.decryptData(ctx, data)
	if err != nil {
		return domain.Data{}, fmt.Errorf("failed decrypt data: %w", err)
	}

	return decryptedData, nil
}
