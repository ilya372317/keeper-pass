package passkeeper

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
)

func (s *Service) All(ctx context.Context) ([]domain.ShortData, error) {
	records, err := s.passClient.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get all records: %w", err)
	}

	return records, nil
}
