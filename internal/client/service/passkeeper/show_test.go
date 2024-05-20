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

func TestService_ShowLoginPass(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: c}

	t.Run("success retrieve case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowLoginPass(gomock.Any(), 1).Times(1).Return(domain.LoginPass{
			URL:      "url",
			Login:    "login",
			Password: "password",
			ID:       1,
		}, nil)

		// Execute.
		got, err := serv.ShowLoginPass(ctx, 1)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "url", got.URL)
		assert.Equal(t, "login", got.Login)
		assert.Equal(t, "password", got.Password)
		assert.Equal(t, 1, got.ID)
	})

	t.Run("failed case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowLoginPass(gomock.Any(), 1).Times(1).Return(domain.LoginPass{},
			fmt.Errorf("internal"))

		// Execute.
		got, err := serv.ShowLoginPass(ctx, 1)

		// Assert.
		require.Nil(t, got)
		require.Error(t, err)
	})
}

func TestService_ShowCreditCard(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: c}

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowCard(gomock.Any(), 1).Times(1).Return(domain.CreditCard{
			BankName:   "bank-name",
			CardNumber: "card-number",
			Exp:        "exp",
			Code:       1,
			ID:         1,
		}, nil)

		// Execute.
		got, err := serv.ShowCreditCard(ctx, 1)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "bank-name", got.BankName)
		assert.Equal(t, "card-number", got.CardNumber)
		assert.Equal(t, "exp", got.Exp)
		assert.Equal(t, 1, got.Code)
		assert.Equal(t, 1, got.ID)
	})

	t.Run("invalid case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowCard(gomock.Any(), 1).Times(1).Return(domain.CreditCard{},
			fmt.Errorf("internal"))

		// Execute.
		got, err := serv.ShowCreditCard(ctx, 1)

		// Assert.
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestService_ShowText(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := passkeeper_mock.NewMockpassClient(ctrl)
	ctx := context.Background()
	serv := Service{passClient: c}

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowText(gomock.Any(), 1).Times(1).Return(domain.Text{
			Info: "info",
			Data: "data",
			ID:   1,
		}, nil)

		// Execute.
		got, err := serv.ShowText(ctx, 1)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "info", got.Info)
		assert.Equal(t, "data", got.Data)
		assert.Equal(t, 1, got.ID)
	})

	t.Run("failed case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowText(gomock.Any(), 1).Times(1).Return(domain.Text{},
			fmt.Errorf("internal"))

		// Execute.
		got, err := serv.ShowText(ctx, 1)

		// Assert.
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestService_ShowBinary(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	c := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: c}

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowBinary(gomock.Any(), 1).Times(1).Return(domain.Binary{
			Info: "info",
			Data: []byte("data"),
			ID:   1,
		}, nil)

		// Execute.
		got, err := serv.ShowBinary(ctx, 1)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "info", got.Info)
		assert.Equal(t, []byte("data"), got.Data)
		assert.Equal(t, 1, got.ID)
	})

	t.Run("failed case", func(t *testing.T) {
		// Prepare.
		c.EXPECT().ShowBinary(gomock.Any(), 1).Times(1).Return(domain.Binary{},
			fmt.Errorf("internal"))

		// Execute.
		got, err := serv.ShowBinary(ctx, 1)

		// Assert.
		require.Error(t, err)
		require.Nil(t, got)
	})
}
