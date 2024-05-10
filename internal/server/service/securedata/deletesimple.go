package securedata

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) DeleteSimple(ctx context.Context, id int64, userID uint) error {
	if err := s.dataStorage.Delete(ctx, []int{int(id)}, userID, domain.KindsCanBeSimpleDeleted); err != nil {
		return fmt.Errorf("failed delete data with id [%d] from storage: %w", id, err)
	}
	return nil
}
