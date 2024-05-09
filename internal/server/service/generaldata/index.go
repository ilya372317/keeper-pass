package generaldata

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) Index(ctx context.Context) ([]domain.GeneralData, error) {
	return []domain.GeneralData{{}}, nil
}
