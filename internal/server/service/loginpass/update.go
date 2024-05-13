package loginpass

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Update(ctx context.Context, d dto.UpdateLoginPassDTO) error {
	data, err := s.dataService.GetAndDecryptData(ctx, d.ID)
	if err != nil {
		return fmt.Errorf("failed get or decrypt data: %w", err)
	}

	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return fmt.Errorf("failed get user from ctx: %w", err)
	}

	if data.UserID != user.ID {
		return fmt.Errorf("failed get data belongs to other user: %w", err)
	}

	if data.Kind != domain.KindLoginPass {
		return fmt.Errorf("wrong data king")
	}

	var updatePayloadDTO dto.LoginPassPayloadDTO
	if err = json.Unmarshal([]byte(data.Payload), &updatePayloadDTO); err != nil {
		return fmt.Errorf("invalid payload in storage: %w", err)
	}
	if d.Metadata != nil {
		metadataContent, err := json.Marshal(d.Metadata)
		if err != nil {
			return fmt.Errorf("failed update metadata: %w", err)
		}
		data.Metadata = string(metadataContent)
	}
	if d.Password != nil {
		updatePayloadDTO.Password = *d.Password
	}
	if d.Login != nil {
		updatePayloadDTO.Login = *d.Login
	}
	savePayloadContent, err := json.Marshal(&updatePayloadDTO)
	if err != nil {
		return fmt.Errorf("failed marshal payload DTO for save")
	}
	data.Payload = string(savePayloadContent)

	if err = s.dataService.EncryptAndUpdateData(ctx, data); err != nil {
		return fmt.Errorf("failed update data in storage: %w", err)
	}

	return nil
}
