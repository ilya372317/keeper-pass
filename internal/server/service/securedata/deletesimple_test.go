package securedata

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	data_mock "github.com/ilya372317/pass-keeper/internal/server/service/securedata/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_DeleteSimple(t *testing.T) {
	ctrl := gomock.NewController(t)
	strg := data_mock.NewMockdataStorage(ctrl)
	serv := Service{dataStorage: strg}
	ctx := context.Background()

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		strg.EXPECT().
			Delete(gomock.Any(), []int{1}, uint(1), domain.KindsCanBeSimpleDeleted).
			Times(1).
			Return(nil)

		// Execute.
		got := serv.DeleteSimple(ctx, 1, 1)

		// Assert.
		require.NoError(t, got)
	})

	t.Run("failed delete in storage", func(t *testing.T) {
		// Prepare.
		strg.EXPECT().
			Delete(gomock.Any(), []int{1}, uint(1), domain.KindsCanBeSimpleDeleted).
			Times(1).
			Return(fmt.Errorf("internal"))

		// Execute.
		got := serv.DeleteSimple(ctx, 1, 1)

		// Assert.
		require.Error(t, got)
	})
}
