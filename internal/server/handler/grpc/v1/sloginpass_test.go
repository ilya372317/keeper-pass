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

func TestHandler_SaveLoginPass(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockloginPassService(ctrl)
	handler := New(v1_mock.NewMockAuthService(ctrl), serv)
	t.Run("success case", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		serv.EXPECT().Save(ctx, gomock.Any()).Times(1).Return(nil)
		arg := &pb.SaveLoginPassRequest{
			Metadata: &pb.LoginPassMetadata{Url: "https://example.com"},
			Login:    "12345",
			Password: "12345",
		}

		// Execute.
		got, err := handler.SaveLoginPass(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("invalid url given", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := &pb.SaveLoginPassRequest{
			Metadata: &pb.LoginPassMetadata{Url: "invalid url"},
			Login:    "12345",
			Password: "12345",
		}

		wantCode := codes.InvalidArgument

		// Execute.
		got, err := handler.SaveLoginPass(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, wantCode, e.Code())
		assert.Nil(t, got)
	})

	t.Run("success case with empty metadata", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		serv.EXPECT().Save(ctx, gomock.Any()).Times(1).Return(nil)
		arg := &pb.SaveLoginPassRequest{
			Metadata: nil,
			Login:    "12345",
			Password: "12345",
		}

		got, err := handler.SaveLoginPass(ctx, arg)
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("invalid password case", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := &pb.SaveLoginPassRequest{
			Metadata: nil,
			Login:    "123",
			Password: "1",
		}
		wantCode := codes.InvalidArgument

		// Execute.
		got, err := handler.SaveLoginPass(ctx, arg)
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, wantCode, e.Code())
		assert.Nil(t, got)
	})

	t.Run("invalid login case", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := &pb.SaveLoginPassRequest{
			Metadata: nil,
			Login:    "1",
			Password: "123",
		}
		wantCode := codes.InvalidArgument

		// Execute.
		got, err := handler.SaveLoginPass(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, wantCode, e.Code())
		assert.Nil(t, got)
	})

	t.Run("internal error on save in service", func(t *testing.T) {
		// Prepare
		ctx := context.Background()
		serv.EXPECT().Save(ctx, gomock.Any()).Times(1).Return(fmt.Errorf("failed save data"))
		arg := &pb.SaveLoginPassRequest{
			Metadata: nil,
			Login:    "123",
			Password: "123",
		}

		// Execute.
		got, err := handler.SaveLoginPass(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, e.Code())
		assert.Nil(t, got)
	})
}

func TestHandler_SaveLoginPass1(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockloginPassService(ctrl)
	handler := New(v1_mock.NewMockAuthService(ctrl), serv)
	t.Run("success update case with all nil value", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		serv.EXPECT().Update(gomock.Any(), dto.UpdateLoginPassDTO{
			ID:       1,
			Metadata: nil,
			Login:    nil,
			Password: nil,
		}).Times(1).Return(nil)
		arg := pb.UpdateLoginPassRequest{
			Id:       1,
			Metadata: nil,
			Login:    nil,
			Password: nil,
		}

		// Execute.
		got, err := handler.UpdateLoginPass(ctx, &arg)

		// Assert.
		require.NoError(t, err)
		assert.NotNil(t, got)
	})
	t.Run("success case with filled args", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		login := stringPtr("123")
		password := stringPtr("123")
		arg := pb.UpdateLoginPassRequest{
			Id:       1,
			Login:    login,
			Password: password,
			Metadata: &pb.LoginPassMetadata{Url: "https://localhost"},
		}
		serv.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(nil)

		// Execute.
		got, err := handler.UpdateLoginPass(ctx, &arg)
		require.NoError(t, err)
		require.NotNil(t, got)
	})
	t.Run("invalid metadata given", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := &pb.UpdateLoginPassRequest{
			Id: 1,
			Metadata: &pb.LoginPassMetadata{
				Url: "invalid-url",
			},
		}

		// Execute.
		got, err := handler.UpdateLoginPass(ctx, arg)

		// Assert.
		require.Error(t, err)
		require.Nil(t, got)
	})
	t.Run("invalid password given", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := &pb.UpdateLoginPassRequest{
			Id:       1,
			Password: stringPtr("1"),
		}

		// Execute.
		got, err := handler.UpdateLoginPass(ctx, arg)
		require.Error(t, err)
		require.Nil(t, got)
	})
	t.Run("invalid login given", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		arg := &pb.UpdateLoginPassRequest{
			Id:    1,
			Login: stringPtr("1"),
		}

		// Execute.
		got, err := handler.UpdateLoginPass(ctx, arg)

		// Assert.
		require.Error(t, err)
		require.Nil(t, got)
	})
	t.Run("failed update in service", func(t *testing.T) {
		// Prepare.
		ctx := context.Background()
		serv.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("internal error"))
		arg := &pb.UpdateLoginPassRequest{
			Id: 1,
		}

		// Execute.
		got, err := handler.UpdateLoginPass(ctx, arg)
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func stringPtr(val string) *string {
	return &val
}
