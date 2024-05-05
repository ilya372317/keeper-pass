package binary

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

type dataService interface {
	EncryptAndSaveData(context.Context, domain.Data) error
	GetAndDecryptData(context.Context, int64) (domain.Data, error)
	EncryptAndUpdateData(context.Context, domain.Data) error
}

type Service struct {
	dataService dataService
}

func New(dataService dataService) *Service {
	return &Service{dataService: dataService}
}
