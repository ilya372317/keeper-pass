package generaldata

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	generaldata_mock "github.com/ilya372317/pass-keeper/internal/server/service/generaldata/mocks"
	"github.com/stretchr/testify/require"
)

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	strg := generaldata_mock.NewMockdataStorage(ctrl)
	serv := Service{dataStorage: strg}
	ctxUser := context.WithValue(context.Background(), domain.CtxUserKey{}, &domain.User{
		ID: 1,
	})

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		strg.EXPECT().DeleteSimple(gomock.Any(), int64(1), uint(1)).Times(1).Return(nil)

		// Execute.
		got := serv.Delete(ctxUser, 1)

		// Assert.
		require.NoError(t, got)
	})

	t.Run("missing user in ctx", func(t *testing.T) {
		// Execute.
		got := serv.Delete(context.Background(), 1)

		// Assert.
		require.Error(t, got)
	})

	t.Run("failed delete in storage", func(t *testing.T) {
		// Prepare.
		strg.EXPECT().DeleteSimple(gomock.Any(), int64(1), uint(1)).Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		got := serv.Delete(ctxUser, 1)

		// Assert.
		require.Error(t, got)
	})
}
