package text

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	text_mock "github.com/ilya372317/pass-keeper/internal/server/service/text/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Show(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := text_mock.NewMockdataService(ctrl)
	serv := Service{dataService: dataServ}
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{
		Email:          "1@gmail.com",
		HashedPassword: "123",
		Salt:           "salt",
		ID:             1,
	})

	t.Run("success show case", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).
			Times(1).Return(domain.Data{
			Payload:            `{"data":"data"}`,
			Metadata:           `{"info":"info"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindText,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		got, err := serv.Show(ctxUser, 1)

		// Assert.
		require.NoError(t, err)
		assert.Equal(t, "data", got.Data)
		assert.Equal(t, "info", got.Metadata.Info)
		assert.Equal(t, int64(1), got.ID)
	})

	t.Run("failed get user from ctx", func(t *testing.T) {
		// Execute.
		_, err := serv.Show(context.Background(), 1)

		// Assert.
		require.Error(t, err)
	})

	t.Run("failed get or decrypt data", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{},
			fmt.Errorf("internal"))

		// Execute.
		_, err := serv.Show(ctxUser, 1)

		// Assert.
		require.Error(t, err)
	})

	t.Run("attempt to get text data belongs to another user", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `{"data":"data"}`,
			Metadata:           `{}`,
			ID:                 1,
			UserID:             2,
			Kind:               domain.KindText,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		_, err := serv.Show(ctxUser, 1)

		// Assert.
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrAccesDenied)
	})

	t.Run("invalid kind", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            "",
			Metadata:           "",
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindBinary,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		_, err := serv.Show(ctxUser, 1)

		// Assert.
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNotSupportedOperation)
	})

	t.Run("invalid payload in storage", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `invalid payload`,
			Metadata:           `{"info":"info"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindText,
			IsPayloadDecrypted: false,
		}, nil)

		// Execute.
		_, err := serv.Show(ctxUser, 1)

		// Assert.
		require.Error(t, err)
	})

	t.Run("invalid metadata in storage", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `{"data":"data"}`,
			Metadata:           `invalid-metadata`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindText,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		_, err := serv.Show(ctxUser, 1)

		// Assert.
		require.Error(t, err)
	})
}
