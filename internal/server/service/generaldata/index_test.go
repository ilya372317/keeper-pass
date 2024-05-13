package generaldata

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	generaldata_mock "github.com/ilya372317/pass-keeper/internal/server/service/generaldata/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_Index(t *testing.T) {
	ctrl := gomock.NewController(t)
	strg := generaldata_mock.NewMockdataStorage(ctrl)
	serv := Service{dataStorage: strg}
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{ID: 1})

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		strg.EXPECT().GetAllEncrypted(gomock.Any(), uint(1)).Times(1).Return([]domain.Data{
			{
				Metadata: `{"url":"info"}`,
				ID:       1,
				UserID:   1,
				Kind:     domain.KindLoginPass,
			},
			{

				Metadata: `{"bank_name":"info"}`,
				ID:       2,
				UserID:   1,
				Kind:     domain.KindCreditCard,
			},
		}, nil)

		// Execute.
		got, err := serv.Index(ctxUser)

		// Assert.
		require.NoError(t, err)
		require.Len(t, got, 2)
	})

	t.Run("missing user in ctx", func(t *testing.T) {
		// Execute.
		got, err := serv.Index(context.Background())

		// Assert.
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("secure data storage return error", func(t *testing.T) {
		// Prepare.
		strg.EXPECT().GetAllEncrypted(gomock.Any(), uint(1)).Times(1).Return(nil,
			fmt.Errorf("internal"))

		// Execute.
		got, err := serv.Index(ctxUser)

		// Assert.
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("failed convert data to specific representation", func(t *testing.T) {
		// Prepare.
		strg.EXPECT().GetAllEncrypted(gomock.Any(), uint(1)).Times(1).Return([]domain.Data{
			{
				Metadata:           "invalid metadata",
				ID:                 1,
				UserID:             1,
				Kind:               domain.KindLoginPass,
				IsPayloadDecrypted: false,
			},
		}, nil)

		// Execute.
		got, err := serv.Index(ctxUser)

		// Assert.
		require.Error(t, err)
		require.Nil(t, got)
	})
}
