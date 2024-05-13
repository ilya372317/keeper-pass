package v1

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	v1_mock "github.com/ilya372317/pass-keeper/internal/server/handler/grpc/v1/mocks"
	pb "github.com/ilya372317/pass-keeper/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_Index(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockgeneralDataService(ctrl)
	handler := Handler{generalDataService: serv}
	ctx := context.Background()
	arg := &pb.IndexRequest{}

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Index(gomock.Any()).Times(1).Return([]domain.GeneralData{
			{
				ID:   1,
				Info: "info",
				Kind: 0,
			},
			{
				ID:   2,
				Info: "info2",
				Kind: 1,
			},
		}, nil)

		// Execute.
		got, err := handler.Index(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Len(t, got.Items, 2)
		firstItem := got.Items[0]
		assert.Equal(t, int64(1), firstItem.Id)
		assert.Equal(t, "info", firstItem.Info)
		assert.Equal(t, int32(0), firstItem.Type)
		secondItem := got.Items[1]
		assert.Equal(t, int64(2), secondItem.Id)
		assert.Equal(t, "info2", secondItem.Info)
		assert.Equal(t, int32(1), secondItem.Type)
	})

	t.Run("failed get data case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Index(gomock.Any()).Times(1).Return(nil, fmt.Errorf("internal"))

		// Execute.
		got, err := handler.Index(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, e.Code())
		assert.Nil(t, got)
	})
}
