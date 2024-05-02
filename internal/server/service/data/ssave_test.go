package data

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/cryptomanager"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	data_mock "github.com/ilya372317/pass-keeper/internal/server/service/data/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctxUser = &domain.User{
	CreatedAT:      time.Now(),
	UpdatedAT:      time.Now(),
	Email:          "email",
	HashedPassword: "pass",
	Salt:           "123",
	ID:             1,
}

var validSecretKey = []byte("1372737473727473")

func TestService_EncryptAndSaveData(t *testing.T) {
	ctrl := gomock.NewController(t)
	strg := data_mock.NewMockdataStorage(ctrl)
	keyr := data_mock.NewMockkeyring(ctrl)
	serv := New(keyr, strg)

	t.Run("success encrypt and save case", func(t *testing.T) {
		// Prepare
		ctx := context.WithValue(context.Background(), domain.CtxUserKey{}, ctxUser)
		keyr.EXPECT().GetGeneralKey(gomock.Any()).Times(1).Return([]byte("1233434312334343"), nil)
		strg.EXPECT().SaveData(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)
		arg := dto.SimpleDataDTO{
			Payload:  `{}`,
			Metadata: `{}`,
			Type:     domain.KindLoginPass,
		}

		// Execute.
		got := serv.EncryptAndSaveData(ctx, arg)

		// Assert.
		require.NoError(t, got)
	})

	t.Run("missing user in ctx", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := dto.SimpleDataDTO{
			Payload:  "{}",
			Metadata: "{}",
			Type:     domain.KindLoginPass,
		}
		// Execute.
		got := serv.EncryptAndSaveData(ctx, arg)
		require.Error(t, got)
	})

	t.Run("failed save data", func(t *testing.T) {
		// Prepare
		ctx := context.WithValue(context.Background(), domain.CtxUserKey{}, ctxUser)
		keyr.EXPECT().GetGeneralKey(gomock.Any()).Times(1).Return([]byte("1233434312334343"), nil)
		strg.EXPECT().SaveData(gomock.Any(), gomock.Any()).
			Times(1).
			Return(fmt.Errorf("failed save"))
		arg := dto.SimpleDataDTO{
			Payload:  `{}`,
			Metadata: `{}`,
			Type:     domain.KindLoginPass,
		}

		// Execute.
		got := serv.EncryptAndSaveData(ctx, arg)

		// Assert.
		require.Error(t, got)
	})
}

func TestService_buildCryptedData(t *testing.T) {
	ctrl := gomock.NewController(t)
	keyr := data_mock.NewMockkeyring(ctrl)
	strg := data_mock.NewMockdataStorage(ctrl)
	serv := New(keyr, strg)

	t.Run("success encryption case", func(t *testing.T) {
		// Prepare.
		ctx := context.WithValue(context.Background(), domain.CtxUserKey{}, ctxUser)
		arg := dto.SimpleDataDTO{
			Payload:  "{}",
			Metadata: "{}",
			Type:     domain.KindLoginPass,
		}
		aesgcm, err := cryptomanager.NewAESGCM(validSecretKey)
		require.NoError(t, err)

		// Execute.
		got, err := serv.buildCryptedData(ctx, arg)
		require.NoError(t, err)
		assert.Len(t, got.Nonce, aesgcm.NonceSize())
	})
}
