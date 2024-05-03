package loginpass

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	loginpass_mock "github.com/ilya372317/pass-keeper/internal/server/service/loginpass/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Show(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := loginpass_mock.NewMockdataService(ctrl)
	serv := New(dataServ)
	userCtx := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{
		CreatedAT:      time.Now(),
		UpdatedAT:      time.Now(),
		Email:          "1@gmail.com",
		HashedPassword: "123",
		Salt:           "salt",
		ID:             1,
	})

	t.Run("success show case", func(t *testing.T) {
		// Prepare.
		arg := 1
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			Payload:            `{"login":"ilya","password":"123"}`,
			Metadata:           `{"url":"https://localhost"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindLoginPass,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		got, err := serv.Show(userCtx, arg)

		// Assert.
		require.NoError(t, err)
		assert.Equal(t, "https://localhost", got.Metadata.URL)
		assert.Equal(t, "ilya", got.Login)
		assert.Equal(t, "123", got.Password)
		assert.Equal(t, 1, got.ID)
	})

	t.Run("empty context", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := 1

		// Execute.
		_, err := serv.Show(ctx, arg)

		// Assert.
		require.Error(t, err)
	})

	t.Run("failed get or decrypt data", func(t *testing.T) {
		// Prepare.
		arg := 1
		dataServ.EXPECT().
			GetAndDecryptData(gomock.Any(), int64(1)).
			Times(1).
			Return(domain.Data{}, fmt.Errorf("failed get data from storage"))

		// Execute.
		_, err := serv.Show(userCtx, arg)

		// Assert.
		require.Error(t, err)
	})

	t.Run("invalid data kind", func(t *testing.T) {
		// Prepare.
		arg := 1
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `{"login":"ilya","password":"123"}`,
			Metadata:           `{"url":"https://localhost"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindFile,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		_, err := serv.Show(userCtx, arg)

		// Assert.
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNotSupportedOperation)
	})

	t.Run("invalid user id", func(t *testing.T) {
		// Prepare.
		arg := 1
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `{"login":"ilya","password":"123"}`,
			Metadata:           `{"url":"https://localhost"}`,
			ID:                 1,
			UserID:             2,
			Kind:               domain.KindLoginPass,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		_, err := serv.Show(userCtx, arg)

		// Assert.
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrAccesDenied)
	})

	t.Run("invalid metadata in storage", func(t *testing.T) {
		// Prepare.
		arg := 1
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `{"login":"ilya","password":"123"}`,
			Metadata:           `invalid-metadata`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindLoginPass,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		_, err := serv.Show(userCtx, arg)

		// Assert.
		require.Error(t, err)
	})

	t.Run("invalid payload in storage", func(t *testing.T) {
		// Prepare.
		arg := 1
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `invalid payload`,
			Metadata:           `{"url":"https://localhost"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindLoginPass,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		_, err := serv.Show(userCtx, arg)

		// Assert.
		require.Error(t, err)
	})
}
