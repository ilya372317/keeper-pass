package passkeeper

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	passkeeper_mock "github.com/ilya372317/pass-keeper/internal/client/service/passkeeper/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_SaveLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: c}
	ctx := context.Background()

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().SaveLogin(gomock.Any(), "login", "pass", "url").Times(1).Return(nil)

		// Execute
		err := serv.SaveLogin(ctx, "login", "pass", "url")

		// Assert.
		require.NoError(t, err)
	})

	t.Run("failed save login in client", func(t *testing.T) {
		// Prepare.
		c.EXPECT().
			SaveLogin(gomock.Any(), "login", "pass", "url").
			Times(1).
			Return(fmt.Errorf("internal"))

		// Execute.
		err := serv.SaveLogin(ctx, "login", "pass", "url")

		// Assert.
		require.Error(t, err)
	})
}
