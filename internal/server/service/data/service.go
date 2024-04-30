package data

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

type keyring interface {
	GetGeneralKey(context.Context) ([]byte, error)
}

type dataStorage interface {
	SaveData(context.Context, domain.Data) error
	GetData(context.Context, int) (*domain.Data, error)
}

type Service struct {
	keyring     keyring
	dataStorage dataStorage
}

func New(keyring keyring, dataStorage dataStorage) *Service {
	return &Service{
		keyring:     keyring,
		dataStorage: dataStorage,
	}
}
