package v1

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	v1_mock "github.com/ilya372317/pass-keeper/internal/server/handler/grpc/v1/mocks"
	pb "github.com/ilya372317/pass-keeper/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_UpdateBinary(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockbinaryService(ctrl)
	handler := Handler{binaryService: serv}
	ctx := context.Background()

	t.Run("success update", func(t *testing.T) {
		// Prepare.
		arg := &pb.UpdateBinaryRequest{
			Id:       1,
			Metadata: &pb.BinaryMetadata{Info: "info"},
			Data:     []byte("321"),
		}
		serv.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(nil)

		// Execute.
		got, err := handler.UpdateBinary(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("failed update", func(t *testing.T) {
		// Prepare.
		arg := &pb.UpdateBinaryRequest{}
		serv.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		got, err := handler.UpdateBinary(ctx, arg)

		// Assert.
		require.Nil(t, got)
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, e.Code())
	})
}
