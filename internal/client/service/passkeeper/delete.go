package passkeeper

import (
	"context"
	"fmt"
	"strconv"
)

func (s *Service) Delete(ctx context.Context, id string) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("id must be integer: %w", err)
	}

	if err = s.passClient.Delete(ctx, intID); err != nil {
		return fmt.Errorf("failed delete: %w", err)
	}

	return nil
}
