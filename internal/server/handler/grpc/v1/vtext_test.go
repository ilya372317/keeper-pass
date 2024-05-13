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

func TestHandler_ShowText(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMocktextService(ctrl)
	handler := Handler{
		textService: serv,
	}
	ctx := context.Background()

	t.Run("success show case", func(t *testing.T) {
		// Prepare.
		arg := &pb.ShowTextRequest{Id: 1}
		serv.EXPECT().Show(gomock.Any(), int64(1)).Times(1).Return(domain.Text{
			Metadata: domain.TextMetadata{
				Info: "info",
			},
			Data: "data",
			ID:   1,
		}, nil)

		// Execute.
		got, err := handler.ShowText(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got.Text)
		require.NotNil(t, got.Text.Metadata)
		assert.Equal(t, "data", got.Text.Data)
		assert.Equal(t, "info", got.Text.Metadata.Info)
		assert.Equal(t, int64(1), got.Text.Id)
	})

	errTests := []struct {
		name      string
		returnErr error
		want      codes.Code
	}{
		{
			name:      "not found",
			returnErr: sql.ErrNoRows,
			want:      codes.NotFound,
		},
		{
			name:      "permission denied",
			returnErr: domain.ErrAccesDenied,
			want:      codes.PermissionDenied,
		},
		{
			name:      "not supported operation",
			returnErr: domain.ErrNotSupportedOperation,
			want:      codes.InvalidArgument,
		},
		{
			name:      "internal",
			returnErr: fmt.Errorf("internal"),
			want:      codes.Internal,
		},
	}

	for _, tt := range errTests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare.
			arg := &pb.ShowTextRequest{Id: 1}
			serv.EXPECT().Show(gomock.Any(), int64(1)).Times(1).Return(domain.Text{}, tt.returnErr)

			// Execute.
			got, err := handler.ShowText(ctx, arg)

			// Assert.
			require.Error(t, err)
			require.Nil(t, got)
			e, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tt.want, e.Code())
		})
	}
}
