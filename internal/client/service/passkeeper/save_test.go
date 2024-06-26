package passkeeper

import (
	"context"
	"fmt"
	"os"
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

func TestService_SaveCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: c}
	ctx := context.Background()

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().SaveCard(gomock.Any(), "number", "exp", 123, "bank name")

		// Execute.
		err := serv.SaveCard(ctx, "number", "exp", "123", "bank name")

		// Assert.
		require.NoError(t, err)
	})

	t.Run("invalid code given", func(t *testing.T) {
		// Execute.
		err := serv.SaveCard(ctx, "number", "exp", "invalid-code", "bank name")

		// Assert.
		require.Error(t, err)
	})

	t.Run("failed save in client", func(t *testing.T) {
		// Prepare.
		c.EXPECT().SaveCard(gomock.Any(), "number", "exp", 123, "bank name").
			Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		err := serv.SaveCard(ctx, "number", "exp", "123", "bank name")

		// Assert.
		require.Error(t, err)
	})
}

func TestService_SaveText(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: c}
	ctx := context.Background()

	t.Run("success save case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().SaveText(gomock.Any(), "info", "data").Times(1).Return(nil)

		// Execute.
		err := serv.SaveText(ctx, "info", "data")

		// Assert.
		require.NoError(t, err)
	})

	t.Run("failed in client", func(t *testing.T) {
		// Prepare.
		c.EXPECT().SaveText(gomock.Any(), "info", "data").Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		err := serv.SaveText(ctx, "info", "data")

		// Assert.
		require.Error(t, err)
	})
}

func TestService_SaveBinary(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: c}
	ctx := context.Background()
	const tempFilePerm = 0600
	validFilePath := "./file.txt"
	err := os.WriteFile(validFilePath, []byte("some data"), tempFilePerm)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.Remove(validFilePath)
		require.NoError(t, err)
	})

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().SaveBinary(gomock.Any(), "info", []byte("some data")).Times(1).Return(nil)

		// Execute.
		err = serv.SaveBinary(ctx, "info", validFilePath)

		// Assert.
		require.NoError(t, err)
	})

	t.Run("given not existing file", func(t *testing.T) {
		// Execute.
		err = serv.SaveBinary(ctx, "info", "./not-existed-file.txt")

		// Assert.
		require.Error(t, err)
	})

	t.Run("error in client", func(t *testing.T) {
		// Prepare.
		c.EXPECT().
			SaveBinary(gomock.Any(), "info", []byte("some data")).
			Times(1).
			Return(fmt.Errorf("internal"))

		// Execute.
		err = serv.SaveBinary(ctx, "info", validFilePath)

		// Assert.
		require.Error(t, err)
	})
}
