package passkeeper

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	passkeeper_mock "github.com/ilya372317/pass-keeper/internal/client/service/passkeeper/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: c}
	ctx := context.Background()

	t.Run("success delete case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().Delete(gomock.Any(), 1).Times(1).Return(nil)

		// Execute.
		err := serv.Delete(ctx, "1")

		// Assert.
		require.NoError(t, err)
	})

	t.Run("invalid id given", func(t *testing.T) {
		// Execute.
		err := serv.Delete(ctx, "invalid-id")

		// Assert.
		require.Error(t, err)
	})

	t.Run("fail in client", func(t *testing.T) {
		// Prepare.
		c.EXPECT().Delete(gomock.Any(), 1).Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		err := serv.Delete(ctx, "1")

		// Assert.
		require.Error(t, err)
	})
}
