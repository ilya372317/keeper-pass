package creditcard

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) Show(ctx context.Context, id int64) (domain.CreditCard, error) {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return domain.CreditCard{}, fmt.Errorf("failed get user from context on credit card show")
	}

	data, err := s.dataService.GetAndDecryptData(ctx, id)
	if err != nil {
		return domain.CreditCard{},
			fmt.Errorf("failed decrypt or get credit card data from storage: %w", err)
	}

	if data.UserID != user.ID {
		return domain.CreditCard{}, domain.ErrAccesDenied
	}

	if data.Kind != domain.KindCreditCard {
		return domain.CreditCard{}, domain.ErrNotSupportedOperation
	}

	var creditCardData domain.CreditCard
	if err = json.Unmarshal([]byte(data.Payload), &creditCardData); err != nil {
		return domain.CreditCard{}, fmt.Errorf("failed parse credit card payload info: %w", err)
	}
	if err = json.Unmarshal([]byte(data.Metadata), &creditCardData.Metadata); err != nil {
		return domain.CreditCard{}, fmt.Errorf("failed parse credit card metadata info: %w", err)
	}
	creditCardData.ID = data.ID

	return creditCardData, nil
}
