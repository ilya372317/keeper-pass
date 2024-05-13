package v1

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	v1_mock "github.com/ilya372317/pass-keeper/internal/server/handler/grpc/v1/mocks"
	pb "github.com/ilya372317/pass-keeper/proto"
	"github.com/stretchr/testify/require"
)

func TestHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockgeneralDataService(ctrl)
	handler := Handler{generalDataService: serv}
	ctx := context.Background()

	t.Run("success delete case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Delete(gomock.Any(), int64(1)).Times(1).Return(nil)
		arg := &pb.DeleteRequest{Id: 1}

		// Execute.
		got, err := handler.Delete(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("failed delete case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Delete(gomock.Any(), int64(1)).Times(1).Return(fmt.Errorf("internal"))
		arg := &pb.DeleteRequest{Id: 1}

		// Execute.
		got, err := handler.Delete(ctx, arg)

		// Assert.
		require.Error(t, err)
		require.Nil(t, got)
	})
}
