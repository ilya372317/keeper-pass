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

func TestHandler_ShowLoginPass(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockloginPassService(ctrl)
	handler := Handler{loginPassService: serv}

	t.Run("success show case", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := &pb.ShowLoginPassRequest{
			Id: 1,
		}
		serv.EXPECT().Show(gomock.Any(), 1).Times(1).Return(domain.LoginPassData{
			Metadata: domain.LoginPassMetadata{
				URL: "https://localhost",
			},
			Login:    "login",
			Password: "pass",
			ID:       1,
		}, nil)

		// Execute.
		got, err := handler.ShowLoginPass(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got.LoginPass)
		require.NotNil(t, got.LoginPass.Metadata)
		assert.Equal(t, "https://localhost", got.LoginPass.Metadata.Url)
		assert.Equal(t, "login", got.LoginPass.Login)
		assert.Equal(t, "pass", got.LoginPass.Password)
	})

	t.Run("data not found case", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := &pb.ShowLoginPassRequest{Id: 1}
		serv.EXPECT().
			Show(gomock.Any(), 1).
			Times(1).
			Return(domain.LoginPassData{}, fmt.Errorf("internal error %w", sql.ErrNoRows))

		// Execute.
		_, err := handler.ShowLoginPass(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, e.Code())
	})

	t.Run("invalid kind case", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := &pb.ShowLoginPassRequest{Id: 1}
		serv.EXPECT().
			Show(gomock.Any(), 1).
			Times(1).
			Return(domain.LoginPassData{}, domain.ErrNotSupportedOperation)

		// Execute.
		_, err := handler.ShowLoginPass(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, e.Code())
	})

	t.Run("internal error", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := &pb.ShowLoginPassRequest{Id: 1}
		serv.EXPECT().Show(gomock.Any(), 1).
			Times(1).
			Return(domain.LoginPassData{}, fmt.Errorf("failed get data"))

		// Execute.
		_, err := handler.ShowLoginPass(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Internal, e.Code())
	})

	t.Run("access denied", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		serv.EXPECT().Show(gomock.Any(), 1).
			Times(1).
			Return(domain.LoginPassData{}, domain.ErrAccesDenied)
		arg := &pb.ShowLoginPassRequest{Id: 1}

		// Execute.
		_, err := handler.ShowLoginPass(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, e.Code())
	})
}
