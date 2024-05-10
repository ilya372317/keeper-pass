package generaldata

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
)

func (s *Service) Index(ctx context.Context) ([]domain.GeneralData, error) {
	user, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
	if !ok {
		return nil, fmt.Errorf("failed get user from context for index data")
	}
	records, err := s.dataStorage.GetAllEncrypted(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed get all data from storage: %w", err)
	}

	generalDataRecords := make([]domain.GeneralData, 0, len(records))

	for i := range records {
		r := &records[i]
		data, err := r.ToDataRepresentation()
		if err != nil {
			return nil, fmt.Errorf("failed covert general data to specific data type")
		}
		generalDataRecords = append(generalDataRecords, domain.GeneralData{
			Info: data.GetInfo(),
			ID:   int64(r.ID),
			Kind: int8(r.Kind),
		})
	}

	return generalDataRecords, nil
}
