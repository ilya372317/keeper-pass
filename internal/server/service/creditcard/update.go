package creditcard

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Update(ctx context.Context, d dto.UpdateCreditCardDTO) error {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return fmt.Errorf("failed get user from ctx for update credit card info")
	}

	data, err := s.dataService.GetAndDecryptData(ctx, int64(d.ID))
	if err != nil {
		return fmt.Errorf("failed get or decrypt data from storage: %w", err)
	}

	if data.UserID != user.ID {
		return domain.ErrAccesDenied
	}

	if data.Kind != domain.KindCreditCard {
		return domain.ErrNotSupportedOperation
	}

	var updatePayloadDTO dto.CreditCardPayload
	if err = json.Unmarshal([]byte(data.Payload), &updatePayloadDTO); err != nil {
		return fmt.Errorf("invalid credit card payload in storage: %w", err)
	}
	if d.Metadata != nil {
		metadataContent, err := json.Marshal(d.Metadata)
		if err != nil {
			return fmt.Errorf("failed update credit card metadata: %w", err)
		}
		data.Metadata = string(metadataContent)
	}
	if d.CardNumber != nil {
		updatePayloadDTO.CardNumber = *d.CardNumber
	}
	if d.Expiration != nil {
		updatePayloadDTO.Expiration = *d.Expiration
	}
	if d.CVV != nil {
		updatePayloadDTO.CVV = int(*d.CVV)
	}
	payloadContent, err := json.Marshal(&updatePayloadDTO)
	if err != nil {
		return fmt.Errorf("failed marshal updated credit card payload: %w", err)
	}
	data.Payload = string(payloadContent)

	if err = s.dataService.EncryptAndUpdateData(ctx, data); err != nil {
		return fmt.Errorf("failed update or encrypt credit card data in storge: %w", err)
	}

	return nil
}
