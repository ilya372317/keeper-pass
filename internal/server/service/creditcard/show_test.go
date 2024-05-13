package creditcard

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	creditcard_mock "github.com/ilya372317/pass-keeper/internal/server/service/creditcard/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Show(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := creditcard_mock.NewMockdataService(ctrl)
	serv := Service{dataService: dataServ}
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{
		Email:          "1@gmail.com",
		HashedPassword: "123",
		Salt:           "salt",
		ID:             1,
	})

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), gomock.Any()).Times(1).Return(domain.Data{
			Payload:            `{"card_number":"123","expiration":"123","cvv":123}`,
			Metadata:           `{"bank_name":"string"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindCreditCard,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		got, err := serv.Show(ctxUser, 1)

		// Assert.
		require.NoError(t, err)
		assert.Equal(t, "123", got.CardNumber)
		assert.Equal(t, "123", got.Expiration)
		assert.Equal(t, 123, got.CVV)
		assert.Equal(t, "string", got.Metadata.BankName)
	})

	t.Run("user not in ctx", func(t *testing.T) {
		// Execute.
		_, err := serv.Show(context.Background(), 1)

		// Assert.
		require.Error(t, err)
	})

	t.Run("invalid user id", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), gomock.Any()).Times(1).Return(domain.Data{
			Payload:            `{"card_number":"123","expiration":"123","cvv":123}`,
			Metadata:           `{"bank_name":"string"}`,
			ID:                 1,
			UserID:             2,
			Kind:               domain.KindCreditCard,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		_, err := serv.Show(ctxUser, 1)

		// Assert.
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAccesDenied)
	})

	t.Run("failed get data from storage", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().
			GetAndDecryptData(gomock.Any(), gomock.Any()).
			Times(1).
			Return(domain.Data{}, fmt.Errorf("internal"))

		// Execute.
		_, err := serv.Show(ctxUser, 1)

		// Assert.
		require.Error(t, err)
	})

	t.Run("invalid kind", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), gomock.Any()).Times(1).Return(domain.Data{
			Payload:            `{"card_number":"123","expiration":"123","cvv":123}`,
			Metadata:           `{"bank_name":"string"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindLoginPass,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		_, err := serv.Show(ctxUser, 1)

		// Assert.
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNotSupportedOperation)
	})
}
