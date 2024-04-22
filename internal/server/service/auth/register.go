package auth

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

// Register method for register new users.
func (s *Service) Register(ctx context.Context, registerDTO dto.RegisterDTO) (*domain.User, error) {
	user := &domain.User{
		Email: registerDTO.Email,
	}
	user.SetHashedPassword(registerDTO.Password)
	if err := user.GenerateSalt(); err != nil {
		return nil, fmt.Errorf("failed generate salt on user reigster: %w", err)
	}

	if err := s.userRepository.SaveUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed save user to database on register: %w", err)
	}

	user, err := s.userRepository.GetUserByEmail(ctx, registerDTO.Email)
	if err != nil {
		return nil, fmt.Errorf("failed get just saved user on register: %w", err)
	}

	return user, nil
}
