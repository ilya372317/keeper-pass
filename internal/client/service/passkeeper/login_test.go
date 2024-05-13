package passkeeper

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	passkeeper_mock "github.com/ilya372317/pass-keeper/internal/client/service/passkeeper/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := passkeeper_mock.NewMockpassClient(ctrl)
	strg := passkeeper_mock.NewMocktokenStorage(ctrl)
	serv := Service{
		passClient:   client,
		tokenStorage: strg,
	}
	ctx := context.Background()

	t.Run("success login case", func(t *testing.T) {
		// Prepare.
		client.EXPECT().Login(gomock.Any(), "email", "pass").Times(1).Return("token", nil)
		strg.EXPECT().SetAccessToken("token").Times(1).Return(nil)

		// Execute.
		err := serv.Login(ctx, "email", "pass")

		// Assert.
		require.NoError(t, err)
	})

	t.Run("failed get token from client", func(t *testing.T) {
		// Prepare.
		client.
			EXPECT().
			Login(gomock.Any(), "email", "pass").
			Times(1).
			Return("", fmt.Errorf("internal"))

		// Execute.
		err := serv.Login(ctx, "email", "pass")

		// Assert.
		require.Error(t, err)
	})

	t.Run("failed save token to storage", func(t *testing.T) {
		// Prepare.
		client.
			EXPECT().
			Login(gomock.Any(), "email", "pass").
			Times(1).
			Return("token", nil)
		strg.EXPECT().SetAccessToken("token").Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		err := serv.Login(ctx, "email", "pass")

		// Assert.
		require.Error(t, err)
	})
}
