package securedata

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) GetAllEncrypted(ctx context.Context) ([]domain.Data, error) {
	dataRecords, err := s.dataStorage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get all data from storage: %w", err)
	}

	return dataRecords, nil
}
