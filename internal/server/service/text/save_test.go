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

func TestService_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := text_mock.NewMockdataService(ctrl)
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{
		Email:          "1@gmail.com",
		HashedPassword: "123",
		Salt:           "salt",
		ID:             1,
	})
	serv := Service{dataService: dataServ}

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		dataServ.EXPECT().EncryptAndSaveData(gomock.Any(), domain.Data{
			Payload:            `{"data":"text"}`,
			Metadata:           `{"info":"test"}`,
			UserID:             1,
			Kind:               domain.KindText,
			IsPayloadDecrypted: true,
		}).Times(1).Return(nil)
		arg := dto.SaveTextDTO{
			Metadata: dto.TextMetadata{
				Info: "test",
			},
			Data: "text",
		}

		// Execute.
		got := serv.Save(ctxUser, arg)

		// Assert.
		require.NoError(t, got)
	})

	t.Run("failed encrypt or save data", func(t *testing.T) {
		// Prepare.
		arg := dto.SaveTextDTO{
			Metadata: dto.TextMetadata{
				Info: "",
			},
			Data: "",
		}
		dataServ.EXPECT().EncryptAndSaveData(gomock.Any(), domain.Data{
			Payload:            `{"data":""}`,
			Metadata:           `{"info":""}`,
			UserID:             1,
			Kind:               domain.KindText,
			IsPayloadDecrypted: true,
		}).Times(1).Return(fmt.Errorf("failed get data"))

		// Execute.
		got := serv.Save(ctxUser, arg)

		// Assert.
		require.Error(t, got)
	})
}
