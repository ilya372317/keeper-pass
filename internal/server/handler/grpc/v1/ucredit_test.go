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

func TestHandler_UpdateCreditCard(t *testing.T) {
	const validCardNumber = "374245455400126"
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockcreditCardService(ctrl)
	handler := Handler{
		creditCardService: serv,
	}
	ctx := context.Background()

	t.Run("success update", func(t *testing.T) {
		// Prepare.
		arg := &pb.UpdateCreditCardRequest{
			Id:         1,
			CardNumber: stringPtr(validCardNumber),
			Expiration: stringPtr("12/20"),
			Cvv:        int32Ptr(123),
			Metadata:   &pb.CreditCardMetadata{BankName: "name"},
		}
		serv.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(nil)

		// Execute.
		got, err := handler.UpdateCreditCard(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("invalid card number", func(t *testing.T) {
		// Prepare.
		arg := &pb.UpdateCreditCardRequest{
			Id:         1,
			CardNumber: stringPtr("invalid-number"),
		}

		// Execute.
		got, err := handler.UpdateCreditCard(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, e.Code())
		assert.Nil(t, got)
	})

	t.Run("invalid expiration", func(t *testing.T) {
		// Prepare.
		arg := &pb.UpdateCreditCardRequest{
			Id:         1,
			Expiration: stringPtr("invalid-expiration"),
		}

		// Execute.
		got, err := handler.UpdateCreditCard(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, e.Code())
		assert.Nil(t, got)
	})

	t.Run("invalid cvv", func(t *testing.T) {
		// Prepare.
		arg := &pb.UpdateCreditCardRequest{
			Id:  1,
			Cvv: int32Ptr(1),
		}

		// Execute.
		got, err := handler.UpdateCreditCard(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, e.Code())
		assert.Nil(t, got)
	})

	t.Run("operation not permitted", func(t *testing.T) {
		// Prepare.
		arg := &pb.UpdateCreditCardRequest{
			Id:         1,
			CardNumber: stringPtr(validCardNumber),
			Expiration: stringPtr("12/20"),
			Cvv:        int32Ptr(123),
			Metadata:   &pb.CreditCardMetadata{BankName: "name"},
		}
		serv.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(domain.ErrAccesDenied)

		// Execute.
		got, err := handler.UpdateCreditCard(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, e.Code())
		assert.Nil(t, got)
	})

	t.Run("operation not permitted", func(t *testing.T) {
		// Prepare.
		arg := &pb.UpdateCreditCardRequest{
			Id:         1,
			CardNumber: stringPtr(validCardNumber),
			Expiration: stringPtr("12/20"),
			Cvv:        int32Ptr(123),
			Metadata:   &pb.CreditCardMetadata{BankName: "name"},
		}
		serv.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		got, err := handler.UpdateCreditCard(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, e.Code())
		assert.Nil(t, got)
	})
}

func int32Ptr(val int32) *int32 {
	return &val
}
