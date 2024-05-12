package passkeeper

import (
	"context"
	"fmt"
)

func (s *Service) SaveLogin(ctx context.Context, login, password, url string) error {
	if err := s.passClient.SaveLogin(ctx, login, password, url); err != nil {
		return fmt.Errorf("failed save login data")
	}

	return nil
}
