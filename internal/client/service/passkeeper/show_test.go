package passkeeper

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/client/domain"
	passkeeper_mock "github.com/ilya372317/pass-keeper/internal/client/service/passkeeper/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Show(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: c}
	ctx := context.Background()

	t.Run("success login pass case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowLoginPass(gomock.Any(), 1).Times(1).Return(domain.LoginPass{
			URL:      "url",
			Login:    "login",
			Password: "pass",
			ID:       1,
		}, nil)

		// Execute.
		res, err := serv.Show(ctx, "1", "login-pass")

		// Assert.
		require.NoError(t, err)
		assert.Equal(t, "ID: 1\nURL: url\nLOGIN: login\nPASSWORD: pass\n", res)
	})

	t.Run("failed get show login pass", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowLoginPass(gomock.Any(), 1).Times(1).Return(domain.LoginPass{}, fmt.Errorf("internal"))

		// Execute.
		_, err := serv.Show(ctx, "1", "login-pass")

		// Assert.
		require.Error(t, err)
	})

	t.Run("success credit card", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowCard(gomock.Any(), 1).Times(1).Return(domain.CreditCard{
			BankName:   "bank",
			CardNumber: "number",
			Exp:        "exp",
			Code:       1,
			ID:         1,
		}, nil)

		// Execute.
		got, err := serv.Show(ctx, "1", "credit-card")

		// Assert.
		require.NoError(t, err)
		assert.Equal(t, "ID : 1\nBANK NAME: bank\nCARD NUMBER: number\nEXP:exp\nCODE:1\n", got)
	})

	t.Run("failed get credit card", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowCard(gomock.Any(), 1).Times(1).Return(domain.CreditCard{}, fmt.Errorf("internal"))

		// Execute.
		_, err := serv.Show(ctx, "1", "credit-card")

		// Assert.
		require.Error(t, err)
	})

	t.Run("success kind text", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowText(gomock.Any(), 1).Times(1).Return(domain.Text{
			Info: "info",
			Data: "data",
			ID:   1,
		}, nil)

		// Execute.
		got, err := serv.Show(ctx, "1", "text")

		// Assert.
		require.NoError(t, err)
		assert.Equal(t, "ID: 1\nINFO: info\nDATA: data\n", got)
	})

	t.Run("failed show text", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowText(gomock.Any(), 1).Times(1).Return(domain.Text{},
			fmt.Errorf("internal"))

		// Execute.
		_, err := serv.Show(ctx, "1", "text")

		// Assert.
		require.Error(t, err)
	})

	t.Run("success binary", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowBinary(gomock.Any(), 1).Times(1).Return(domain.Binary{
			Info: "info",
			Data: []byte("data"),
			ID:   1,
		}, nil)

		// Execute.
		got, err := serv.Show(ctx, "1", "binary")

		// Assert.
		require.NoError(t, err)
		assert.Equal(t, "ID: 1\nINFO: info\nDATA: data\n", got)
	})

	t.Run("failed show binary", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowBinary(gomock.Any(), 1).Times(1).Return(domain.Binary{},
			fmt.Errorf("internal"))

		// Execute.
		_, err := serv.Show(ctx, "1", "binary")

		// Assert.
		require.Error(t, err)
	})

	t.Run("given alias is invalid", func(t *testing.T) {
		// Execute.
		_, err := serv.Show(ctx, "1", "invalid-kind")

		// Assert.
		require.Error(t, err)
	})

	t.Run("invalid id given", func(t *testing.T) {
		// Execute.
		_, err := serv.Show(ctx, "invalid-id", "binary")

		// Assert.
		require.Error(t, err)
	})
}
