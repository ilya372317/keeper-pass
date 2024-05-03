package loginpass

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	loginpass_mock "github.com/ilya372317/pass-keeper/internal/server/service/loginpass/mocks"
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

func TestService_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	dataServ := loginpass_mock.NewMockdataService(ctrl)
	serv := New(dataServ)
	ctxWithUser := context.WithValue(context.Background(), domain.CtxUserKey{}, ctxUser)

	t.Run("success save case", func(t *testing.T) {
		// Setup.
		ctx := ctxWithUser
		dataServ.EXPECT().EncryptAndSaveData(ctx, domain.Data{
			Payload:            `{"login":"123","password":"123"}`,
			Metadata:           `{"url":"https://localhost"}`,
			Kind:               domain.KindLoginPass,
			UserID:             1,
			IsPayloadDecrypted: true,
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
		ctx := ctxWithUser
		dataServ.EXPECT().EncryptAndSaveData(ctx, domain.Data{
			Payload:            `{"login":"123","password":"123"}`,
			Metadata:           `{}`,
			Kind:               domain.KindLoginPass,
			UserID:             1,
			IsPayloadDecrypted: true,
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

	t.Run("missing user in context", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := dto.SaveLoginPassDTO{
			Metadata: dto.LoginPassMetadata{},
			Login:    "123",
			Password: "123",
		}
		// Execute.
		got := serv.Save(ctx, arg)

		// Assert.
		require.Error(t, got)
	})
}
