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

func TestHandler_UpdateText(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMocktextService(ctrl)
	handler := Handler{textService: serv}
	ctx := context.Background()

	t.Run("success save case", func(t *testing.T) {
		// Prepare.
		arg := &pb.UpdateTextRequest{
			Id:       1,
			Metadata: &pb.TextMetadata{Info: "info"},
			Data:     stringPtr("data"),
		}
		serv.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(nil)

		// Execute.
		got, err := handler.UpdateText(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("failed update in service case", func(t *testing.T) {
		// Prepare.
		arg := &pb.UpdateTextRequest{
			Id: 1,
		}
		serv.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		got, err := handler.UpdateText(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, e.Code())
		assert.Nil(t, got)
	})
}
