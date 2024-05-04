package creditcard

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	creditcard_mock "github.com/ilya372317/pass-keeper/internal/server/service/creditcard/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := creditcard_mock.NewMockdataService(ctrl)
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{
		Email:          "1@gmail.com",
		HashedPassword: "123",
		Salt:           "salt",
		ID:             1,
	})
	serv := Service{dataService: dataServ}

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateCreditCardDTO{
			Metadata:   &dto.CreditCardMetadata{BankName: "string"},
			CardNumber: stringPtr("123"),
			Expiration: stringPtr("123"),
			CVV:        int32Ptr(123),
			ID:         1,
		}
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			ID:                 1,
			Payload:            `{}`,
			Metadata:           `{}`,
			UserID:             1,
			Kind:               domain.KindCreditCard,
			IsPayloadDecrypted: true,
		}, nil)
		dataServ.EXPECT().EncryptAndUpdateData(gomock.Any(), domain.Data{
			Payload:            `{"card_number":"123","expiration":"123","cvv":123}`,
			Metadata:           `{"bank_name":"string"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindCreditCard,
			IsPayloadDecrypted: true,
		}).Times(1).Return(nil)

		// Execute.
		got := serv.Update(ctxUser, arg)

		// Assert.
		require.NoError(t, got)
	})

	t.Run("invalid payload in storage", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateCreditCardDTO{
			Metadata:   &dto.CreditCardMetadata{BankName: "string"},
			CardNumber: stringPtr("123"),
			Expiration: stringPtr("123"),
			CVV:        int32Ptr(123),
			ID:         1,
		}
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			ID:                 1,
			Payload:            `invalid-payload`,
			Metadata:           `{}`,
			UserID:             1,
			Kind:               domain.KindCreditCard,
			IsPayloadDecrypted: true,
		}, nil)

		// Execute.
		got := serv.Update(ctxUser, arg)

		// Assert.
		require.Error(t, got)
	})

	t.Run("failed get or decrypt data from storage", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateCreditCardDTO{
			Metadata:   &dto.CreditCardMetadata{BankName: "string"},
			CardNumber: stringPtr("123"),
			Expiration: stringPtr("123"),
			CVV:        int32Ptr(123),
			ID:         1,
		}
		dataServ.
			EXPECT().
			GetAndDecryptData(gomock.Any(), int64(1)).
			Times(1).
			Return(domain.Data{}, fmt.Errorf("failed get data"))

		// Execute.
		got := serv.Update(ctxUser, arg)

		// Assert.
		require.Error(t, got)
	})

	t.Run("failed update or encrypt data", func(t *testing.T) {
		// Prepare.
		arg := dto.UpdateCreditCardDTO{
			Metadata:   &dto.CreditCardMetadata{BankName: "string"},
			CardNumber: stringPtr("123"),
			Expiration: stringPtr("123"),
			CVV:        int32Ptr(123),
			ID:         1,
		}
		dataServ.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(domain.Data{
			ID:                 1,
			Payload:            `{}`,
			Metadata:           `{}`,
			UserID:             1,
			Kind:               domain.KindCreditCard,
			IsPayloadDecrypted: true,
		}, nil)
		dataServ.EXPECT().EncryptAndUpdateData(gomock.Any(), domain.Data{
			Payload:            `{"card_number":"123","expiration":"123","cvv":123}`,
			Metadata:           `{"bank_name":"string"}`,
			ID:                 1,
			UserID:             1,
			Kind:               domain.KindCreditCard,
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

func int32Ptr(val int32) *int32 {
	return &val
}
