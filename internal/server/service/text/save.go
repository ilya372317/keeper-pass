package text

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Save(ctx context.Context, d dto.SaveTextDTO) error {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return fmt.Errorf("missing user in ctx for make text update")
	}

	textPayload := dto.TextPayload{Data: d.Data}
	textPayloadContent, err := json.Marshal(&textPayload)
	if err != nil {
		return fmt.Errorf("failed marshal text payload: %w", err)
	}

	textMetadataContent, err := json.Marshal(&d.Metadata)
	if err != nil {
		return fmt.Errorf("failed marshal text metadata: %w", err)
	}

	data := domain.Data{
		Payload:            string(textPayloadContent),
		Metadata:           string(textMetadataContent),
		UserID:             user.ID,
		Kind:               domain.KindText,
		IsPayloadDecrypted: true,
	}

	if err = s.dataService.EncryptAndSaveData(ctx, data); err != nil {
		return fmt.Errorf("failed encrypt or save text data: %w", err)
	}

	return nil
}
