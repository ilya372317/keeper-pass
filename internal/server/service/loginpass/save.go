package loginpass

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Save(ctx context.Context, d dto.SaveLoginPassDTO) error {
	metadataString, err := json.Marshal(&d.Metadata)
	if err != nil {
		return fmt.Errorf("failed marshal metadata: %w", err)
	}
	payloadString, err := json.Marshal(&dto.LoginPassPayloadDTO{
		Login:    d.Login,
		Password: d.Password,
	})
	if err != nil {
		return fmt.Errorf("failed marshal login pass payload: %w", err)
	}
	sd := dto.SimpleDataDTO{
		Payload:  string(payloadString),
		Metadata: string(metadataString),
		Type:     domain.KindLoginPass,
	}

	if err = s.dataService.EncryptAndSaveData(ctx, sd); err != nil {
		return fmt.Errorf("failed ecnrypt and save data: %w", err)
	}

	return nil
}
