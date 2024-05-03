package data

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) EncryptAndSaveData(ctx context.Context, d dto.SaveSimpleDataDTO) error {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return fmt.Errorf("failed get user from context")
	}

	newData := domain.Data{
		Payload:            d.Payload,
		Metadata:           d.Metadata,
		UserID:             user.ID,
		Kind:               d.Type,
		IsPayloadDecrypted: true,
	}
	encryptedData, err := s.encryptData(ctx, newData)
	if err != nil {
		return fmt.Errorf("failed encrypt data: %w", err)
	}

	if err = s.dataStorage.SaveData(ctx, encryptedData); err != nil {
		return fmt.Errorf("failed save data to storage: %w", err)
	}

	return nil
}
