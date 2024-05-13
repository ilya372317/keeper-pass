package keyring

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/service/cryptomanager"
	keyring_mock "github.com/ilya372317/pass-keeper/internal/server/service/keyring/mocks"
	"github.com/stretchr/testify/require"
)

func TestKeyring_GetGeneralKeyWithKeyInStorage(t *testing.T) {
	defer func() {
		cryptoKey = nil
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyStorage := keyring_mock.NewMockkeyStorage(ctrl)
	masterKey := []byte("123423gjduvydhsj")
	k := New(masterKey, mockKeyStorage)

	ctx := context.Background()

	secretKey := []byte("secretKeydjvhduz")
	aesgcm, err := cryptomanager.NewAESGCM(masterKey)
	require.NoError(t, err)
	nonce, err := cryptomanager.GenerateRandom(aesgcm.NonceSize())
	require.NoError(t, err)
	secretCipherKey := cryptomanager.Encrypt(aesgcm, secretKey, nonce)

	mockKeyStorage.
		EXPECT().
		GetKey(ctx).
		Return(
			&domain.Keys{Key: hex.EncodeToString(secretCipherKey), Nonce: hex.EncodeToString(nonce)}, nil,
		).
		Times(1)
	mockKeyStorage.EXPECT().SaveKey(ctx, gomock.Any()).Return(nil).Times(0)

	_, err = k.GetGeneralKey(ctx)
	require.NoError(t, err)

	_, err = k.GetGeneralKey(ctx)
	require.NoError(t, err)
}

func TestKeyring_GetGeneralKey_ErrorRetrievingKey(t *testing.T) {
	defer func() {
		cryptoKey = nil
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyStorage := keyring_mock.NewMockkeyStorage(ctrl)
	k := New([]byte("masterKeydjgkvjd"), mockKeyStorage)

	ctx := context.Background()

	mockKeyStorage.EXPECT().GetKey(ctx).Return(nil, errors.New("failed to retrieve key")).Times(1)

	_, err := k.GetGeneralKey(ctx)
	require.Error(t, err)
}

func TestKeyring_GetGeneralKeyInvalidMasterKeySize(t *testing.T) {
	defer func() {
		cryptoKey = nil
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyStorage := keyring_mock.NewMockkeyStorage(ctrl)
	k := New([]byte("invalidkey"), mockKeyStorage)

	ctx := context.Background()

	_, err := k.GetGeneralKey(ctx)
	require.Error(t, err)
}

func TestKeyring_InitGeneralKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyStorage := keyring_mock.NewMockkeyStorage(ctrl)
	masterKey := []byte("masterKey12345678901234567890123")
	k := New(masterKey, mockKeyStorage)

	ctx := context.Background()

	t.Run("success with existing key", func(t *testing.T) {
		defer func() { cryptoKey = nil }()

		secretKey := []byte("secretKey1234fdj")
		aesgcm, err := cryptomanager.NewAESGCM(masterKey)
		require.NoError(t, err)
		nonce, err := cryptomanager.GenerateRandom(aesgcm.NonceSize())
		require.NoError(t, err)
		secretCipherKey := cryptomanager.Encrypt(aesgcm, secretKey, nonce)

		mockKeyStorage.EXPECT().
			GetKey(ctx).
			Return(&domain.Keys{Key: hex.EncodeToString(secretCipherKey), Nonce: hex.EncodeToString(nonce)}, nil).
			Times(1)

		err = k.InitGeneralKey(ctx)

		require.NoError(t, err)
		require.Equal(t, secretKey, cryptoKey)
	})

	t.Run("success with new key generation", func(t *testing.T) {
		defer func() { cryptoKey = nil }()

		mockKeyStorage.EXPECT().
			GetKey(ctx).
			Return(nil, sql.ErrNoRows).
			Times(1)

		mockKeyStorage.EXPECT().
			SaveKey(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, keys *domain.Keys) error {
				decodedKey, _ := hex.DecodeString(keys.Key)
				decodedNonce, _ := hex.DecodeString(keys.Nonce)
				aesgcm, _ := cryptomanager.NewAESGCM(masterKey)
				decryptedKey, err := cryptomanager.Decrypt(aesgcm, decodedKey, decodedNonce)
				require.NoError(t, err)
				require.Len(t, decryptedKey, generalKeySize)
				return nil
			}).
			Times(1)

		err := k.InitGeneralKey(ctx)

		require.NoError(t, err)
		require.Len(t, cryptoKey, generalKeySize)
	})

	t.Run("failure on key storage error", func(t *testing.T) {
		defer func() { cryptoKey = nil }()

		mockKeyStorage.EXPECT().
			GetKey(ctx).
			Return(nil, errors.New("storage error")).
			Times(1)

		err := k.InitGeneralKey(ctx)

		require.Error(t, err)
	})
}
