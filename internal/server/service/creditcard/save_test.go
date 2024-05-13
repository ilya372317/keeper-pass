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

func TestService_Save(t *testing.T) {
	defaultArg := dto.SaveCreditCardDTO{
		Metadata: dto.CreditCardMetadata{
			BankName: "bank name",
		},
		CardNumber: "123",
		Expiration: "exp",
		CVV:        123,
	}
	defaultData := domain.Data{
		Payload:            `{"card_number":"123","expiration":"exp","cvv":123}`,
		Metadata:           `{"bank_name":"bank name"}`,
		UserID:             1,
		Kind:               domain.KindCreditCard,
		IsPayloadDecrypted: true,
	}
	ctrl := gomock.NewController(t)
	dataStorage := creditcard_mock.NewMockdataService(ctrl)
	serv := Service{dataService: dataStorage}
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{
		Email:          "1@gmail.com",
		HashedPassword: "123",
		Salt:           "salt",
		ID:             1,
	})

	t.Run("success save", func(t *testing.T) {
		// Prepare.
		dataStorage.EXPECT().EncryptAndSaveData(gomock.Any(), defaultData).Times(1).Return(nil)

		// Execute.
		got := serv.Save(ctxUser, defaultArg)

		// Assert.
		require.NoError(t, got)
	})

	t.Run("missing user in context", func(t *testing.T) {
		// Prepare.
		arg := dto.SaveCreditCardDTO{}
		// Execute.
		got := serv.Save(context.Background(), arg)

		// Assert.
		require.Error(t, got)
	})

	t.Run("failed save or encrypt data", func(t *testing.T) {
		// Prepare.
		dataStorage.
			EXPECT().
			EncryptAndSaveData(gomock.Any(), defaultData).
			Times(1).
			Return(fmt.Errorf("failed save"))

		// Execute.
		got := serv.Save(ctxUser, defaultArg)

		// Assert.
		require.Error(t, got)
	})
}
