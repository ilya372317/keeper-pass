package text

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Update(ctx context.Context, d dto.UpdateTextDTO) error {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return fmt.Errorf("failed get user from context for update text")
	}

	data, err := s.dataService.GetAndDecryptData(ctx, d.ID)
	if err != nil {
		return fmt.Errorf("failed get or decrypt data for update text: %w", err)
	}

	if data.Kind != domain.KindText {
		return domain.ErrNotSupportedOperation
	}

	if data.UserID != user.ID {
		return domain.ErrAccesDenied
	}

	var (
		updatePayload  dto.TextPayload
		updateMetadata dto.TextMetadata
	)

	if err = json.Unmarshal([]byte(data.Payload), &updatePayload); err != nil {
		return fmt.Errorf("invalid text payload in storage: %w", err)
	}

	if err = json.Unmarshal([]byte(data.Metadata), &updateMetadata); err != nil {
		return fmt.Errorf("invalid text metadata in storage: %w", err)
	}

	if d.Data != nil {
		updatePayload.Data = *d.Data
	}

	if d.Metadata != nil {
		updateMetadata.Info = d.Metadata.Info
	}

	var (
		metadataContent []byte
		payloadContent  []byte
	)

	if payloadContent, err = json.Marshal(&updatePayload); err != nil {
		return fmt.Errorf("failed marshal text payload")
	}
	if metadataContent, err = json.Marshal(&updateMetadata); err != nil {
		return fmt.Errorf("failed marshal text metadata")
	}

	data.Payload = string(payloadContent)
	data.Metadata = string(metadataContent)

	if err = s.dataService.EncryptAndUpdateData(ctx, data); err != nil {
		return fmt.Errorf("failed update or encrypt text data")
	}

	return nil
}
