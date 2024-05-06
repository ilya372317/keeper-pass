package text

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	text_mock "github.com/ilya372317/pass-keeper/internal/server/service/text/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := text_mock.NewMockdataService(ctrl)
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{
		Email:          "1@gmail.com",
		HashedPassword: "123",
		Salt:           "salt",
		ID:             1,
	})
	serv := Service{dataService: dataServ}

	t.Run("success update case", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateTextDTO{
			Metadata: &dto.TextMetadata{Info: "test"},
			Data:     stringPtr("123"),
			ID:       1,
		}
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `{"data":"data"}`,
			Metadata:           `{"info":"info"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindText,
			IsPayloadDecrypted: true,
		}, nil)
		dataServ.EXPECT().EncryptAndUpdateData(gomock.Any(), domain.Data{
			Payload:            `{"data":"123"}`,
			Metadata:           `{"info":"test"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindText,
			IsPayloadDecrypted: true,
		}).Times(1).Return(nil)

		// Execute.
		got := serv.Update(ctxUser, arg)

		// Assert.
		require.NoError(t, got)
	})

	t.Run("missing user in ctx", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateTextDTO{}

		// Execute.
		got := serv.Update(context.Background(), arg)

		// Assert.
		require.Error(t, got)
	})

	t.Run("invalid kind", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateTextDTO{
			ID: 1,
		}
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `{"data":"data"}`,
			Metadata:           `{"info":"info"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindBinary,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		err := serv.Update(ctxUser, arg)

		// Assert.
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNotSupportedOperation)
	})

	t.Run("failed get or decrypt data", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateTextDTO{
			ID: 1,
		}
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(
			domain.Data{},
			fmt.Errorf("internal"),
		)

		// Execute.
		err := serv.Update(ctxUser, arg)

		// Assert.
		require.Error(t, err)
	})

	t.Run("success update case", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateTextDTO{
			Metadata: &dto.TextMetadata{Info: "test"},
			Data:     stringPtr("123"),
			ID:       1,
		}
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			Payload:            `{"data":"data"}`,
			Metadata:           `{"info":"info"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindText,
			IsPayloadDecrypted: true,
		}, nil)
		dataServ.EXPECT().EncryptAndUpdateData(gomock.Any(), domain.Data{
			Payload:            `{"data":"123"}`,
			Metadata:           `{"info":"test"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindText,
			IsPayloadDecrypted: true,
		}).Times(1).Return(fmt.Errorf("failed update data"))

		// Execute.
		got := serv.Update(ctxUser, arg)

		// Assert.
		require.Error(t, got)
	})
}

func stringPtr(val string) *string {
	return &val
}
