package v1

import (
	"context"
	"database/sql"
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

func TestHandler_ShowBinary(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockbinaryService(ctrl)
	handler := Handler{binaryService: serv}
	ctx := context.Background()

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		arg := &pb.ShowBinaryRequest{Id: 1}
		serv.EXPECT().Show(gomock.Any(), int64(1)).Times(1).Return(domain.Binary{
			Metadata: domain.BinaryMetadata{
				Info: "info",
			},
			Data: []byte("data"),
			ID:   1,
		}, nil)

		// Execute.
		got, err := handler.ShowBinary(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
		require.NotNil(t, got.Binary)
		require.NotNil(t, got.Binary.Metadata)
		assert.Equal(t, []byte("data"), got.Binary.Data)
		assert.Equal(t, "info", got.Binary.Metadata.Info)
		assert.Equal(t, int64(1), got.Binary.Id)
	})

	errTests := []struct {
		name string
		err  error
		want codes.Code
	}{
		{
			name: "not found",
			err:  sql.ErrNoRows,
			want: codes.NotFound,
		},
		{
			name: "permission denied",
			err:  domain.ErrAccesDenied,
			want: codes.PermissionDenied,
		},
		{
			name: "operation not supported",
			err:  domain.ErrNotSupportedOperation,
			want: codes.InvalidArgument,
		},
		{
			name: "internal error",
			err:  fmt.Errorf("internal"),
			want: codes.Internal,
		},
	}
	for _, tt := range errTests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare.
			arg := &pb.ShowBinaryRequest{Id: 1}
			serv.EXPECT().Show(gomock.Any(), gomock.Any()).Times(1).Return(domain.Binary{}, tt.err)

			// Execute.
			got, err := handler.ShowBinary(ctx, arg)

			// Assert.
			require.Nil(t, got)
			require.Error(t, err)
			e, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tt.want, e.Code())
		})
	}
}
