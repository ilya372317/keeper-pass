package loginpass

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

type dataService interface {
	EncryptAndSaveData(context.Context, dto.SimpleDataDTO) error
}

type Service struct {
	dataService dataService
}

func New(dataService dataService) *Service {
	return &Service{
		dataService: dataService,
	}
}
