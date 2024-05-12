package passkeeper

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
)

type passClient interface {
	Login(context.Context, string, string) (string, error)
	All(ctx context.Context) ([]domain.ShortData, error)
	Register(context.Context, string, string) error
	SaveLogin(context.Context, string, string, string) error
	SaveCard(context.Context, string, string, int, string) error
	SaveText(context.Context, string, string) error
	SaveBinary(context.Context, string, []byte) error
	Delete(context.Context, int) error
}

type tokenStorage interface {
	SetAccessToken(token string) error
}

type Service struct {
	passClient   passClient
	tokenStorage tokenStorage
}

func New(client passClient, tokenStorage tokenStorage) *Service {
	return &Service{
		passClient:   client,
		tokenStorage: tokenStorage,
	}
}
