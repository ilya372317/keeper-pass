package securedata

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) EncryptAndUpdateData(ctx context.Context, data domain.Data) error {
	encryptedData, err := s.encryptData(ctx, data)
	if err != nil {
		return fmt.Errorf("failed ecnrypt data before update: %w", err)
	}

	if err = s.dataStorage.UpdateByID(ctx, data.ID, encryptedData); err != nil {
		return fmt.Errorf("failed update data in storage")
	}

	return nil
}
