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

func TestHandler_SaveText(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMocktextService(ctrl)
	handler := Handler{
		textService: serv,
	}
	ctx := context.Background()

	t.Run("success save case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Save(gomock.Any(), dto.SaveTextDTO{
			Metadata: dto.TextMetadata{
				Info: "info",
			},
			Data: "test data",
		}).Times(1).Return(nil)
		arg := &pb.SaveTextRequest{
			Metadata: &pb.TextMetadata{Info: "info"},
			Data:     "test data",
		}

		// Execute.
		got, err := handler.SaveText(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("failed save data", func(t *testing.T) {
		// Prepare.
		arg := &pb.SaveTextRequest{}
		serv.EXPECT().Save(gomock.Any(), dto.SaveTextDTO{}).Times(1).Return(fmt.Errorf("failed get data"))

		// Execute.
		got, err := handler.SaveText(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, e.Code())
		assert.Nil(t, got)
	})
}
