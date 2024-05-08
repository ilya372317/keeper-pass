package securedata

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) EncryptAndSaveData(ctx context.Context, data domain.Data) error {
	encryptedData, err := s.encryptData(ctx, data)
	if err != nil {
		return fmt.Errorf("failed encrypt data: %w", err)
	}

	if err = s.dataStorage.SaveData(ctx, encryptedData); err != nil {
		return fmt.Errorf("failed save data to storage: %w", err)
	}

	return nil
}
