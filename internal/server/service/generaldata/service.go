package generaldata

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

type dataStorage interface {
	GetAllEncrypted(context.Context, uint) ([]domain.Data, error)
	DeleteSimple(context.Context, int64, uint) error
}

type Service struct {
	dataStorage dataStorage
}

func New(dataStorage dataStorage) *Service {
	return &Service{dataStorage: dataStorage}
}
