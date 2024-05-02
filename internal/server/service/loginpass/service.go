package loginpass

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

type dataService interface {
	EncryptAndSaveData(context.Context, dto.SaveSimpleDataDTO) error
	GetAndDecryptData(context.Context, int64) (*domain.Data, error)
	EncryptAndUpdateData(context.Context, dto.UpdateSimpleDataDTO) error
}

type Service struct {
	dataService dataService
}

func New(dataService dataService) *Service {
	return &Service{
		dataService: dataService,
	}
}
