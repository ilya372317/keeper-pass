package binary

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) Show(ctx context.Context, id int64) (domain.Binary, error) {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return domain.Binary{}, fmt.Errorf("failed get user from ctx for show binary")
	}

	data, err := s.dataService.GetAndDecryptData(ctx, id)
	if err != nil {
		return domain.Binary{}, fmt.Errorf("failed get or decrypt data: %w", err)
	}

	if data.UserID != user.ID {
		return domain.Binary{}, domain.ErrAccesDenied
	}

	if data.Kind != domain.KindBinary {
		return domain.Binary{}, domain.ErrNotSupportedOperation
	}

	var (
		result domain.Binary
	)

	if err = json.Unmarshal([]byte(data.Payload), &result); err != nil {
		return domain.Binary{}, fmt.Errorf("failed unmarshal binary payload: %w", err)
	}

	if err = json.Unmarshal([]byte(data.Metadata), &result.Metadata); err != nil {
		return domain.Binary{}, fmt.Errorf("faild unmarshal binary metadata: %w", err)
	}
	result.ID = int64(data.ID)

	return result, nil
}
