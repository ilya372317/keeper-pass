package generaldata

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

type dataStorage interface {
	GetAll(context.Context) ([]domain.Data, error)
}

type Service struct {
	dataStorage dataStorage
}

func New(dataStorage dataStorage) *Service {
	return &Service{dataStorage: dataStorage}
}
