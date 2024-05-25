package passkeeper

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/client/domain"
	passkeeper_mock "github.com/ilya372317/pass-keeper/internal/client/service/passkeeper/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_UpdateLoginPass(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	validArg := &domain.LoginPass{
		URL:      "url",
		Login:    "login",
		Password: "password",
		ID:       1,
	}
	serv := Service{passClient: c}

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().UpdateLoginPass(gomock.Any(), validArg).Times(1).Return(nil)

		// Execute.
		err := serv.UpdateLoginPass(ctx, validArg)

		// Assert.
		require.NoError(t, err)
	})

	t.Run("failed case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().UpdateLoginPass(gomock.Any(), validArg).Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		err := serv.UpdateLoginPass(ctx, validArg)

		// Assert.
		require.Error(t, err)
	})
}
