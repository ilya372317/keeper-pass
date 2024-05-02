package loginpass

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	loginpass_mock "github.com/ilya372317/pass-keeper/internal/server/service/loginpass/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := loginpass_mock.NewMockdataService(ctrl)
	serv := New(dataServ)

	t.Run("success save case", func(t *testing.T) {
		// Setup.
		ctx := context.Background()
		dataServ.EXPECT().EncryptAndSaveData(ctx, dto.SaveSimpleDataDTO{
			Payload:  `{"login":"123","password":"123"}`,
			Metadata: `{"url":"https://localhost"}`,
			Type:     domain.KindLoginPass,
		}).
			Times(1).
			Return(nil)
		arg := dto.SaveLoginPassDTO{
			Metadata: dto.LoginPassMetadata{
				URL: "https://localhost",
			},
			Login:    "123",
			Password: "123",
		}

		// Execute.
		got := serv.Save(ctx, arg)

		// Assert.
		require.NoError(t, got)
	})

	t.Run("failed encrypt and save data case", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		dataServ.EXPECT().EncryptAndSaveData(ctx, dto.SaveSimpleDataDTO{
			Payload:  `{"login":"123","password":"123"}`,
			Metadata: `{}`,
			Type:     domain.KindLoginPass,
		}).
			Times(1).
			Return(fmt.Errorf("internal"))
		arg := dto.SaveLoginPassDTO{
			Metadata: dto.LoginPassMetadata{},
			Login:    "123",
			Password: "123",
		}

		// Execute.
		got := serv.Save(ctx, arg)
		require.Error(t, got)
	})
}
