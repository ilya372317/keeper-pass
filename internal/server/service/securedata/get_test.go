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

func TestService_GetAndDecryptData(t *testing.T) {
	ctrl := gomock.NewController(t)
	keyr := data_mock.NewMockkeyring(ctrl)
	strg := data_mock.NewMockdataStorage(ctrl)
	serv := New(keyr, strg)

	t.Run("success get", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		var arg int64 = 1
		wantPayload := "123"
		wantMetadata := "321"
		keyr.EXPECT().GetGeneralKey(gomock.Any()).Times(2).Return(validSecretKey, nil)
		encryptedDataInStorage, err := serv.encryptData(ctx, domain.Data{
			Payload:            wantPayload,
			Metadata:           wantMetadata,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindLoginPass,
			IsPayloadDecrypted: true,
		})
		require.NoError(t, err)
		strg.EXPECT().GetDataByID(ctx, 1).Times(1).Return(encryptedDataInStorage, nil)

		// Execute.
		got, err := serv.GetAndDecryptData(ctx, arg)

		// Assert.
		require.NoError(t, err)
		assert.Equal(t, wantPayload, got.Payload)
		assert.Equal(t, wantMetadata, got.Metadata)
	})

	t.Run("failed get data from storage", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		var arg int64 = 1
		strg.EXPECT().
			GetDataByID(gomock.Any(), 1).
			Times(1).
			Return(domain.Data{}, fmt.Errorf("internal error"))

		// Execute.
		_, err := serv.GetAndDecryptData(ctx, arg)

		// Assert.
		require.Error(t, err)
	})

	t.Run("failed decrypt data in storage", func(t *testing.T) {
		// Prepare.
		keyr.EXPECT().GetGeneralKey(gomock.Any()).Times(2).Return(validSecretKey, nil)
		ctx := context.Background()
		var arg int64 = 1
		encryptedDataInStorage, err := serv.encryptData(ctx, domain.Data{})
		require.NoError(t, err)
		encryptedDataInStorage.Payload = "invalid-payload"

		strg.EXPECT().GetDataByID(gomock.Any(), 1).Times(1).Return(encryptedDataInStorage, nil)

		// Execute.
		_, err = serv.GetAndDecryptData(ctx, arg)
		// Assert.
		require.Error(t, err)
	})
}
