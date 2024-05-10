package securedata

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	data_mock "github.com/ilya372317/pass-keeper/internal/server/service/securedata/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_GetAllEncrypted(t *testing.T) {
	ctrl := gomock.NewController(t)
	strg := data_mock.NewMockdataStorage(ctrl)
	serv := Service{dataStorage: strg}
	ctx := context.Background()

	t.Run("success get all case", func(t *testing.T) {
		// Prepare.
		strg.EXPECT().GetAll(gomock.Any(), uint(1)).Times(1).Return([]domain.Data{
			{}, {}, {},
		}, nil)

		// Execute.
		got, err := serv.GetAllEncrypted(ctx, 1)

		// Assert.
		require.NoError(t, err)
		assert.Len(t, got, 3)
	})

	t.Run("failed get data from storage", func(t *testing.T) {
		// Prepare.
		strg.EXPECT().GetAll(gomock.Any(), uint(1)).Times(1).Return(nil, fmt.Errorf("internal"))

		// Execute.
		got, err := serv.GetAllEncrypted(ctx, 1)

		// Assert.
		require.Error(t, err)
		assert.Nil(t, got)
	})
}
