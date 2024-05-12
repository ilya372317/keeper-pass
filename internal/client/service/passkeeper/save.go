package passkeeper

import (
	"context"
	"fmt"
	"strconv"
)

func (s *Service) SaveLogin(ctx context.Context, login, password, url string) error {
	if err := s.passClient.SaveLogin(ctx, login, password, url); err != nil {
		return fmt.Errorf("failed save login data: %w", err)
	}

	return nil
}

func (s *Service) SaveCard(ctx context.Context, number, exp, code, bankName string) error {
	intCode, err := strconv.Atoi(code)
	if err != nil {
		return fmt.Errorf("code must be integer: %w", err)
	}

	if err = s.passClient.SaveCard(ctx, number, exp, intCode, bankName); err != nil {
		return fmt.Errorf("failed save card data: %w", err)
	}

	return nil
}

func (s *Service) SaveText(ctx context.Context, info, data string) error {
	if err := s.passClient.SaveText(ctx, info, data); err != nil {
		return fmt.Errorf("failed save text info: %w", err)
	}

	return nil
}
