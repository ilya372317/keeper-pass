package text

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) Show(ctx context.Context, id int64) (domain.Text, error) {
	return domain.Text{}, nil
}
