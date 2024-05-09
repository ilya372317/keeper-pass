package loginpass

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

func (s *Service) Show(ctx context.Context, id int) (domain.LoginPass, error) {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return domain.LoginPass{}, fmt.Errorf("failed get user from context")
	}

	data, err := s.dataService.GetAndDecryptData(ctx, int64(id))
	if err != nil {
		return domain.LoginPass{}, fmt.Errorf("failed get or decrypt data: %w", err)
	}

	if data.Kind != domain.KindLoginPass {
		return domain.LoginPass{}, domain.ErrNotSupportedOperation
	}

	if data.UserID != user.ID {
		return domain.LoginPass{}, domain.ErrAccesDenied
	}

	var loginPassMetadata dto.LoginPassMetadata
	if err = json.Unmarshal([]byte(data.Metadata), &loginPassMetadata); err != nil {
		return domain.LoginPass{}, fmt.Errorf("invalid metadata in storage: %w", err)
	}

	var loginPassPayload dto.LoginPassPayloadDTO
	if err = json.Unmarshal([]byte(data.Payload), &loginPassPayload); err != nil {
		return domain.LoginPass{}, fmt.Errorf("invalid payload in storage: %w", err)
	}

	return domain.LoginPass{
		Metadata: domain.LoginPassMetadata{
			URL: loginPassMetadata.URL,
		},
		Login:    loginPassPayload.Login,
		Password: loginPassPayload.Password,
		ID:       data.ID,
	}, nil
}
