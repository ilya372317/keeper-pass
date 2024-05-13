package binary

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Update(ctx context.Context, d dto.UpdateBinaryDTO) error {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return fmt.Errorf("failed get user from ctx for update binary data")
	}

	data, err := s.dataService.GetAndDecryptData(ctx, d.ID)
	if err != nil {
		return fmt.Errorf("failed get or decrypt binary data: %w", err)
	}

	if data.UserID != user.ID {
		return domain.ErrAccesDenied
	}

	if data.Kind != domain.KindBinary {
		return domain.ErrNotSupportedOperation
	}

	var (
		binaryPayload  dto.BinaryPayload
		binaryMetadata dto.BinaryMetadata
	)

	if err = json.Unmarshal([]byte(data.Payload), &binaryPayload); err != nil {
		return fmt.Errorf("failed unmarshal binary payload for update: %w", err)
	}

	if err = json.Unmarshal([]byte(data.Metadata), &binaryMetadata); err != nil {
		return fmt.Errorf("failed unmarshal binary metadata for update: %w", err)
	}

	if d.Data != nil {
		binaryPayload.Data = *d.Data
	}
	if d.Metadata != nil {
		binaryMetadata.Info = d.Metadata.Info
	}

	var (
		payloadContent  []byte
		payloadMetadata []byte
	)

	payloadContent, err = json.Marshal(&binaryPayload)
	if err != nil {
		return fmt.Errorf("failed marshal binary payload for update: %w", err)
	}
	payloadMetadata, err = json.Marshal(&binaryMetadata)
	if err != nil {
		return fmt.Errorf("failed marshal binary metadata for update: %w", err)
	}

	data.Payload = string(payloadContent)
	data.Metadata = string(payloadMetadata)

	if err = s.dataService.EncryptAndUpdateData(ctx, data); err != nil {
		return fmt.Errorf("failed encrypt or update binary data: %w", err)
	}

	return nil
}
