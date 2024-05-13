package binary

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Save(ctx context.Context, d dto.SaveBinaryDTO) error {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return fmt.Errorf("failed get user from ctx for save binary data")
	}

	binaryPayload := dto.BinaryPayload{Data: d.Data}
	binaryMetadata := dto.BinaryMetadata{Info: d.Metadata.Info}

	binaryPayloadContent, err := json.Marshal(&binaryPayload)
	if err != nil {
		return fmt.Errorf("failed marshal binary payload: %w", err)
	}

	binaryMetadataContent, err := json.Marshal(&binaryMetadata)
	if err != nil {
		return fmt.Errorf("failed marshal binary metadata: %w", err)
	}

	data := domain.Data{
		Payload:            string(binaryPayloadContent),
		Metadata:           string(binaryMetadataContent),
		UserID:             user.ID,
		Kind:               domain.KindBinary,
		IsPayloadDecrypted: true,
	}

	if err = s.dataService.EncryptAndSaveData(ctx, data); err != nil {
		return fmt.Errorf("failed encrypt or save data: %w", err)
	}

	return nil
}
