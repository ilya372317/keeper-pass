package passkeeper

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
)

func (s *Service) UpdateLoginPass(ctx context.Context, lp *domain.LoginPass) error {
	if err := s.passClient.UpdateLoginPass(ctx, lp); err != nil {
		return fmt.Errorf("failed update login pass with id #%d: %w", lp.ID, err)
	}

	return nil
}
