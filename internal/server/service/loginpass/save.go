package loginpass

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Save(ctx context.Context, d dto.SaveLoginPassDTO) error {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return fmt.Errorf("failed get user from context")
	}

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
	data := domain.Data{
		Payload:            string(payloadString),
		Metadata:           string(metadataString),
		UserID:             user.ID,
		Kind:               domain.KindLoginPass,
		IsPayloadDecrypted: true,
	}

	if err = s.dataService.EncryptAndSaveData(ctx, data); err != nil {
		return fmt.Errorf("failed ecnrypt and save data: %w", err)
	}

	return nil
}
