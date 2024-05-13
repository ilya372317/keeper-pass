package text

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) Show(ctx context.Context, id int64) (domain.Text, error) {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return domain.Text{}, fmt.Errorf("failed get user from context for show text")
	}

	data, err := s.dataService.GetAndDecryptData(ctx, id)
	if err != nil {
		return domain.Text{}, fmt.Errorf("failed get or decrypt data for show text: %w", err)
	}

	if user.ID != data.UserID {
		return domain.Text{}, domain.ErrAccesDenied
	}

	if data.Kind != domain.KindText {
		return domain.Text{}, domain.ErrNotSupportedOperation
	}

	var (
		payloadText  domain.Text
		metadataText domain.TextMetadata
	)

	if err = json.Unmarshal([]byte(data.Payload), &payloadText); err != nil {
		return domain.Text{}, fmt.Errorf("failed unmarshal text payload: %w", err)
	}

	if err = json.Unmarshal([]byte(data.Metadata), &metadataText); err != nil {
		return domain.Text{}, fmt.Errorf("failed unmarshal text metadata: %w", err)
	}

	return domain.Text{
		Metadata: domain.TextMetadata{
			Info: metadataText.Info,
		},
		Data: payloadText.Data,
		ID:   int64(data.ID),
	}, nil
}
