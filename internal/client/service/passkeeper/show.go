package passkeeper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
)

func (s *Service) Show(ctx context.Context, id, kindAlias string) (string, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return "", fmt.Errorf("id must be integer for show: %w", err)
	}

	kind, ok := domain.AliasToKind[kindAlias]
	if !ok {
		return "", fmt.Errorf("given alias is invalid")
	}

	var data domain.ShowAble

	switch kind {
	case domain.KindLoginPass:
		data, err = s.passClient.ShowLoginPass(ctx, intID)
		if err != nil {
			return "", fmt.Errorf("failed get login pass info: %w", err)
		}
	case domain.KindCreditCard:
		data, err = s.passClient.ShowCard(ctx, intID)
		if err != nil {
			return "", fmt.Errorf("failed get credit card info: %w", err)
		}
	case domain.KindText:
		data, err = s.passClient.ShowText(ctx, intID)
		if err != nil {
			return "", fmt.Errorf("failed get text info: %w", err)
		}
	case domain.KindBinary:
		data, err = s.passClient.ShowBinary(ctx, intID)
		if err != nil {
			return "", fmt.Errorf("failed get binary info: %w", err)
		}
	}

	if data == nil {
		return "", fmt.Errorf("kind not implementing for show")
	}

	return data.ToString(), nil
}

// ShowLoginPass retrieve login pass from client.
func (s *Service) ShowLoginPass(ctx context.Context, id int) (*domain.LoginPass, error) {
	lp, err := s.passClient.ShowLoginPass(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed get login pass from client: %w", err)
	}

	return &lp, nil
}

// ShowCreditCard retrieve credit card from client.
func (s *Service) ShowCreditCard(ctx context.Context, id int) (*domain.CreditCard, error) {
	cc, err := s.passClient.ShowCard(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed get credit card data from client: %w", err)
	}

	return &cc, nil
}

func (s *Service) ShowText(ctx context.Context, id int) (*domain.Text, error) {
	text, err := s.passClient.ShowText(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed get text from client: %w", err)
	}

	return &text, nil
}

func (s *Service) ShowBinary(ctx context.Context, id int) (*domain.Binary, error) {
	binary, err := s.passClient.ShowBinary(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed get binary from client: %w", err)
	}

	return &binary, nil
}
