package securedata

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

// GetAllEncrypted retrieves all data from storage without decrypting their payload.
// It may be used if the client code does not need to work with the payload.
// For example, if the client needs to list metadata or other data information, this method could be useful.
func (s *Service) GetAllEncrypted(ctx context.Context) ([]domain.Data, error) {
	dataRecords, err := s.dataStorage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get all data from storage: %w", err)
	}

	return dataRecords, nil
}
