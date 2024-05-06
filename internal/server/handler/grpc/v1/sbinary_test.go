package v1

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	v1_mock "github.com/ilya372317/pass-keeper/internal/server/handler/grpc/v1/mocks"
	pb "github.com/ilya372317/pass-keeper/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_SaveBinary(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockbinaryService(ctrl)
	handler := Handler{binaryService: serv}
	ctx := context.Background()

	t.Run("success save case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Save(gomock.Any(), dto.SaveBinaryDTO{
			Metadata: dto.TextMetadata{
				Info: "info",
			},
			Data: []byte("data"),
		}).Times(1).Return(nil)
		arg := &pb.SaveBinaryRequest{
			Metadata: &pb.BinaryMetadata{Info: "info"},
			Data:     []byte("data"),
		}

		// Execute.
		got, err := handler.SaveBinary(ctx, arg)

		// Assert.
		require.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("failed save in service", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("internal"))
		arg := &pb.SaveBinaryRequest{}

		// Execute.
		got, err := handler.SaveBinary(ctx, arg)

		// Assert.
		require.Error(t, err)
		require.Nil(t, got)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, e.Code())
	})
}
