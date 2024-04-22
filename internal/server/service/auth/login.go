package auth

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Login(ctx context.Context, dto dto.LoginDTO) (string, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return "", fmt.Errorf("failed get user from storage for login: %w", err)
	}

	if passCorrect := user.IsPasswordCorrect(dto.Password); !passCorrect {
		return "", domain.ErrInvalidPassword
	}

	tokenString, err := s.tokenManager.Generate(&user)
	if err != nil {
		return "", fmt.Errorf("failed generate auth token: %w", err)
	}

	return tokenString, nil
}
