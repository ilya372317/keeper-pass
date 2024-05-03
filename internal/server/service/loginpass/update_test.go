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

func TestService_Update(t *testing.T) {
	const userID = 1
	var defaultDataInStorage = &domain.Data{
		Payload:  `{"login":"password"}`,
		Metadata: `{"url":"https://localhost"}`,
		ID:       1,
		UserID:   userID,
		Kind:     domain.KindLoginPass,
	}
	ctrl := gomock.NewController(t)
	dataStorage := loginpass_mock.NewMockdataService(ctrl)
	serv := New(dataStorage)
	ctxWithUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{
		CreatedAT:      time.Now(),
		UpdatedAT:      time.Now(),
		Email:          "1@gmail.com",
		HashedPassword: "password",
		Salt:           "salt",
		ID:             userID,
	})

	t.Run("success update case", func(t *testing.T) {
		// Prepare
		arg := dto.UpdateLoginPassDTO{
			Metadata: &dto.LoginPassMetadata{URL: "https://localhost:80"},
			Login:    stringPtr("login"),
			Password: stringPtr("123"),
			ID:       1,
		}
		dataStorage.EXPECT().GetAndDecryptData(gomock.Any(), gomock.Any()).Times(1).Return(defaultDataInStorage, nil)
		dataStorage.EXPECT().EncryptAndUpdateData(gomock.Any(), dto.UpdateSimpleDataDTO{
			Payload:  `{"login":"login","password":"123"}`,
			Metadata: `{"url":"https://localhost:80"}`,
			Type:     domain.KindLoginPass,
			ID:       1,
		}).Times(1).Return(nil)

		// Execute.
		got := serv.Update(ctxWithUser, arg)

		// Assert.
		require.NoError(t, got)
	})
	t.Run("without user in ctx", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		dataStorage.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(defaultDataInStorage, nil)
		arg := dto.UpdateLoginPassDTO{
			ID: 1,
		}

		// Execute.
		got := serv.Update(ctx, arg)

		// Assert.
		require.Error(t, got)
	})
	t.Run("failed get data from storage case", func(t *testing.T) {
		// Prepare.
		ctx := ctxWithUser
		arg := dto.UpdateLoginPassDTO{
			ID: 1,
		}
		dataStorage.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(nil, fmt.Errorf("internal error"))

		// Execute.
		got := serv.Update(ctx, arg)

		// Assert.
		require.Error(t, got)
	})

	t.Run("attempt to get data belongs to other person", func(t *testing.T) {
		// Prepare.
		ctx := ctxWithUser
		arg := dto.UpdateLoginPassDTO{
			ID: 1,
		}
		dataStorage.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(&domain.Data{
			ID:     1,
			UserID: 2,
			Kind:   domain.KindLoginPass,
		}, nil)

		// Execute.
		got := serv.Update(ctx, arg)

		// Assert.
		require.Error(t, got)
	})

	t.Run("invalid kind update requested", func(t *testing.T) {
		// Prepare.
		ctx := ctxWithUser
		arg := dto.UpdateLoginPassDTO{
			ID: 1,
		}
		dataStorage.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(&domain.Data{
			ID:     1,
			UserID: userID,
			Kind:   domain.KindFile,
		}, nil)

		// Execute.
		got := serv.Update(ctx, arg)

		// Assert.
		require.Error(t, got)
	})

	t.Run("failed update data is storage", func(t *testing.T) {
		// Prepare.
		ctx := ctxWithUser
		arg := dto.UpdateLoginPassDTO{
			ID: 1,
		}
		dataStorage.EXPECT().GetAndDecryptData(gomock.Any(), int64(1)).Times(1).Return(&domain.Data{
			Payload:  "{}",
			Metadata: "{}",
			ID:       1,
			UserID:   userID,
			Kind:     domain.KindLoginPass,
		}, nil)
		dataStorage.EXPECT().EncryptAndUpdateData(gomock.Any(), dto.UpdateSimpleDataDTO{
			Payload:  "{}",
			Metadata: "{}",
			Type:     domain.KindLoginPass,
			ID:       1,
		}).Times(1).Return(fmt.Errorf("internal error"))

		// Execute.
		got := serv.Update(ctx, arg)

		// Assert.
		require.Error(t, got)
	})
}

func stringPtr(val string) *string {
	return &val
}
