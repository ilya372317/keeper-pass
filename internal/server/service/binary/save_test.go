package binary

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	binary_mock "github.com/ilya372317/pass-keeper/internal/server/service/binary/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := binary_mock.NewMockdataService(ctrl)
	serv := Service{dataService: dataServ}
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{
		Email:          "1@gmail.com",
		HashedPassword: "123",
		Salt:           "salt",
		ID:             1,
	})

	t.Run("success save case", func(t *testing.T) {
		// Prepare.
		arg := dto.SaveBinaryDTO{
			Metadata: dto.TextMetadata{
				Info: "info",
			},
			Data: []byte("data"),
		}
		dataServ.EXPECT().EncryptAndSaveData(gomock.Any(), domain.Data{
			Payload:            `{"data":"ZGF0YQ=="}`,
			Metadata:           `{"info":"info"}`,
			UserID:             1,
			Kind:               domain.KindBinary,
			IsPayloadDecrypted: true,
		}).Times(1).Return(nil)

		// Execute.
		err := serv.Save(ctxUser, arg)

		// Assert.
		require.NoError(t, err)
	})

	t.Run("user missing in ctx", func(t *testing.T) {
		// Prepare.
		arg := dto.SaveBinaryDTO{}

		// Execute.
		err := serv.Save(context.Background(), arg)

		// Assert.
		require.Error(t, err)
	})

	t.Run("failed save", func(t *testing.T) {
		// Prepare.
		arg := dto.SaveBinaryDTO{}
		dataServ.EXPECT().EncryptAndSaveData(gomock.Any(), gomock.Any()).Times(1).
			Return(fmt.Errorf("internal"))

		// Execute.
		err := serv.Save(ctxUser, arg)

		// Assert.
		require.Error(t, err)
	})
}
