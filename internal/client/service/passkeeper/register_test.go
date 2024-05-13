package passkeeper

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	passkeeper_mock "github.com/ilya372317/pass-keeper/internal/client/service/passkeeper/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: c}
	ctx := context.Background()

	t.Run("success register case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().Register(gomock.Any(), "email", "pass").Times(1).Return(nil)

		// Execute.
		err := serv.Register(ctx, "email", "pass")

		// Assert.
		require.NoError(t, err)
	})

	t.Run("failed in client", func(t *testing.T) {
		// Prepare
		c.EXPECT().Register(ctx, "email", "pass").Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		err := serv.Register(ctx, "email", "pass")

		// Assert.
		require.Error(t, err)
	})
}
