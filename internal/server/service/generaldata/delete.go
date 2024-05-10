package generaldata

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) Delete(ctx context.Context, id int64) error {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return fmt.Errorf("failed get user from ctx for delete data")
	}

	if err := s.dataStorage.DeleteSimple(ctx, id, user.ID); err != nil {
		return fmt.Errorf("failed delete simple data: %w", err)
	}

	return nil
}
