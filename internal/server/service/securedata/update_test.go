package securedata

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	data_mock "github.com/ilya372317/pass-keeper/internal/server/service/securedata/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_EncryptAndUpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	keyr := data_mock.NewMockkeyring(ctrl)
	strg := data_mock.NewMockdataStorage(ctrl)
	serv := New(keyr, strg)

	t.Run("success update case", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := domain.Data{
			ID: 1,
		}
		keyr.EXPECT().GetGeneralKey(gomock.Any()).Times(1).Return(validSecretKey, nil)
		strg.EXPECT().UpdateByID(gomock.Any(), 1, gomock.Any()).Times(1).Return(nil)

		// Execute.
		got := serv.EncryptAndUpdateData(ctx, arg)

		// Assert.
		require.NoError(t, got)
	})

	t.Run("failed encrypt data", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := domain.Data{
			ID: 1,
		}
		keyr.EXPECT().GetGeneralKey(gomock.Any()).Times(1).Return(nil, fmt.Errorf("internal error"))

		// Execute.
		got := serv.EncryptAndUpdateData(ctx, arg)

		// Assert.
		require.Error(t, got)
	})

	t.Run("failed update data in stogage", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := domain.Data{
			ID: 1,
		}
		keyr.EXPECT().GetGeneralKey(gomock.Any()).Times(1).Return(validSecretKey, nil)
		strg.EXPECT().UpdateByID(ctx, 1, gomock.Any()).Times(1).Return(fmt.Errorf("internal error"))

		// Execute.
		got := serv.EncryptAndUpdateData(ctx, arg)

		// Assert.
		require.Error(t, got)
	})
}
