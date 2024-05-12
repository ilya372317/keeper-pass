package passkeeper

import (
	"context"
	"fmt"
)

func (s *Service) Register(ctx context.Context, email string, password string) error {
	if err := s.passClient.Register(ctx, email, password); err != nil {
		return fmt.Errorf("failed register user: %w", err)
	}
	return nil
}
