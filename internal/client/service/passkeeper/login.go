package passkeeper

import (
	"context"
	"fmt"
)

func (s *Service) Login(ctx context.Context, email, password string) error {
	token, err := s.passClient.Login(ctx, email, password)
	if err != nil {
		return fmt.Errorf("failed get token by client: %w", err)
	}
	if err = s.tokenStorage.SetAccessToken(token); err != nil {
		return fmt.Errorf("failed store access token: %w", err)
	}

	return nil
}
