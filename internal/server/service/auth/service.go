package auth

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

type TokenManager interface {
	Generate(user *domain.User) (string, error)
}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	SaveUser(ctx context.Context, user *domain.User) error
	HasUser(ctx context.Context, email string) (bool, error)
}

type Service struct {
	tokenManager   TokenManager
	userRepository UserRepository
}

func NewAuthService(tokenManager TokenManager, userRepository UserRepository) *Service {
	return &Service{
		tokenManager:   tokenManager,
		userRepository: userRepository,
	}
}
