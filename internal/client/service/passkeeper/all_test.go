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

func TestService_All(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := passkeeper_mock.NewMockpassClient(ctrl)
	serv := Service{passClient: client}
	ctx := context.Background()

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		client.EXPECT().All(gomock.Any()).Times(1).Return([]domain.ShortData{{}, {}}, nil)

		// Execute.
		got, err := serv.All(ctx)

		// Assert.
		require.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("internal error", func(t *testing.T) {
		// Prepare.
		client.EXPECT().All(gomock.Any()).Times(1).Return(nil, fmt.Errorf("internal"))

		// Execute.
		got, err := serv.All(ctx)

		// Assert.
		require.Error(t, err)
		assert.Nil(t, got)
	})
}
