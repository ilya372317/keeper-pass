package data

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/service/cryptomanager"
	data_mock "github.com/ilya372317/pass-keeper/internal/server/service/data/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_encryptData(t *testing.T) {
	ctrl := gomock.NewController(t)
	keyr := data_mock.NewMockkeyring(ctrl)
	strg := data_mock.NewMockdataStorage(ctrl)
	serv := New(keyr, strg)

	t.Run("success encryption case", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		keyr.EXPECT().GetGeneralKey(gomock.Any()).Times(1).Return(validSecretKey, nil)
		payload := `{"login":"password"}`
		arg := domain.Data{
			Payload:            payload,
			Metadata:           "{}",
			UserID:             1,
			Kind:               domain.KindLoginPass,
			IsPayloadDecrypted: true,
		}
		aesgcm, err := cryptomanager.NewAESGCM(validSecretKey)
		require.NoError(t, err)

		// Execute.
		got, err := serv.encryptData(ctx, arg)
		require.NoError(t, err)

		// Assert.
		gotPayloadNonce, err := hex.DecodeString(got.PayloadNonce)
		require.NoError(t, err)
		gotKey, err := hex.DecodeString(got.CryptoKey)
		require.NoError(t, err)
		gotDataKeyNonce, err := hex.DecodeString(got.CryptoKeyNonce)
		require.NoError(t, err)
		gotDecryptedDataKey, err := aesgcm.Open(nil, gotDataKeyNonce, gotKey, nil)
		require.NoError(t, err)
		dataAesgcm, err := cryptomanager.NewAESGCM(gotDecryptedDataKey)
		require.NoError(t, err)
		gotPayload, err := hex.DecodeString(got.Payload)
		require.NoError(t, err)
		decryptedPayload, err := dataAesgcm.Open(nil, gotPayloadNonce, gotPayload, nil)
		require.NoError(t, err)

		assert.Len(t, gotDecryptedDataKey, domain.CryptoKeyLength)
		assert.Len(t, gotPayloadNonce, aesgcm.NonceSize())
		assert.Len(t, gotPayloadNonce, dataAesgcm.NonceSize())
		assert.Equal(t, payload, string(decryptedPayload))
	})

	t.Run("invalid secret key in keyring", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		keyr.EXPECT().GetGeneralKey(gomock.Any()).Times(1).Return([]byte("invalid-key"), nil)
		arg := domain.Data{
			Payload:            "123",
			Metadata:           "123",
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindLoginPass,
			IsPayloadDecrypted: true,
		}

		// Execute.
		_, err := serv.encryptData(ctx, arg)

		// Assert.
		require.Error(t, err)
	})
}

func TestService_decryptData(t *testing.T) {
	ctrl := gomock.NewController(t)
	keyr := data_mock.NewMockkeyring(ctrl)
	serv := New(keyr, data_mock.NewMockdataStorage(ctrl))
	dataKeyAesgcm, err := cryptomanager.NewAESGCM(validSecretKey)
	require.NoError(t, err)
	dataKeyNonce, err := cryptomanager.GenerateRandom(dataKeyAesgcm.NonceSize())
	require.NoError(t, err)
	dataKey, err := cryptomanager.GenerateRandom(domain.CryptoKeyLength)
	require.NoError(t, err)
	cryptedDataKey := cryptomanager.Encrypt(dataKeyAesgcm, dataKey, dataKeyNonce)
	require.NoError(t, err)
	payload := `{"login":"ilya","password":"123"}`
	payloadAesgcm, err := cryptomanager.NewAESGCM(dataKey)
	require.NoError(t, err)
	payloadNonce, err := cryptomanager.GenerateRandom(payloadAesgcm.NonceSize())
	require.NoError(t, err)
	cryptedPayload := cryptomanager.Encrypt(payloadAesgcm, []byte(payload), payloadNonce)

	t.Run("success decryption case", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := domain.Data{
			Payload:            hex.EncodeToString(cryptedPayload),
			PayloadNonce:       hex.EncodeToString(payloadNonce),
			CryptoKeyNonce:     hex.EncodeToString(dataKeyNonce),
			CryptoKey:          hex.EncodeToString(cryptedDataKey),
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindLoginPass,
			IsPayloadDecrypted: false,
		}
		keyr.EXPECT().GetGeneralKey(gomock.Any()).Times(1).Return(validSecretKey, nil)

		// Execution.
		data, err := serv.decryptData(ctx, arg)

		// Assert.
		require.NoError(t, err)
		assert.Equal(t, payload, data.Payload)
	})
}
