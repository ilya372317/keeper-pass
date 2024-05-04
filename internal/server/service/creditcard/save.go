package creditcard

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Save(ctx context.Context, d dto.SaveCreditCardDTO) error {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return fmt.Errorf("user not found in ctx")
	}

	payloadData, err := json.Marshal(&d)
	if err != nil {
		return fmt.Errorf("failed marshal credit card payload for save")
	}

	metadataData, err := json.Marshal(&d.Metadata)
	if err != nil {
		return fmt.Errorf("failed marshal credit card metadata for save")
	}

	data := domain.Data{
		Payload:            string(payloadData),
		Metadata:           string(metadataData),
		UserID:             user.ID,
		Kind:               domain.KindCreditCard,
		IsPayloadDecrypted: true,
	}

	if err = s.dataService.EncryptAndSaveData(ctx, data); err != nil {
		return fmt.Errorf("failed ecnrypt or save credit card info to storage: %w", err)
	}

	return nil
}
