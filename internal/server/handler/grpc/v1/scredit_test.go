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

func TestHandler_SaveCreditCard(t *testing.T) {
	const validCardNumber = "374245455400126"
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockcreditCardService(ctrl)
	handler := Handler{
		creditCardService: serv,
	}
	ctx := context.Background()

	t.Run("invalid credit card number given", func(t *testing.T) {
		// Prepare.
		arg := &pb.SaveCreditCardRequest{
			Metadata:   nil,
			CardNumber: "card-number",
			Expiration: "01/23",
			Cvv:        123,
		}

		// Execute
		got, err := handler.SaveCreditCard(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, e.Code())
		assert.Nil(t, got)
	})

	t.Run("invalid expiration date", func(t *testing.T) {
		// Prepare.
		arg := &pb.SaveCreditCardRequest{
			Metadata:   &pb.CreditCardMetadata{BankName: "test-name"},
			CardNumber: validCardNumber,
			Expiration: "13/20",
			Cvv:        123,
		}

		// Execute.
		got, err := handler.SaveCreditCard(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, e.Code())
		assert.Nil(t, got)
	})

	t.Run("to low secret code", func(t *testing.T) {
		// Prepare.
		arg := &pb.SaveCreditCardRequest{
			Metadata:   &pb.CreditCardMetadata{BankName: "test-name"},
			CardNumber: validCardNumber,
			Expiration: "12/20",
			Cvv:        10,
		}

		// Execute.
		got, err := handler.SaveCreditCard(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, e.Code())
		assert.Nil(t, got)
	})

	t.Run("to high secret code", func(t *testing.T) {
		// Prepare.
		arg := &pb.SaveCreditCardRequest{
			Metadata:   &pb.CreditCardMetadata{BankName: "test-name"},
			CardNumber: validCardNumber,
			Expiration: "12/20",
			Cvv:        99999,
		}

		// Execute.
		got, err := handler.SaveCreditCard(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, e.Code())
		assert.Nil(t, got)
	})

	t.Run("failed save data in service", func(t *testing.T) {
		// Prepare.
		arg := &pb.SaveCreditCardRequest{
			Metadata:   &pb.CreditCardMetadata{BankName: "bank name"},
			CardNumber: validCardNumber,
			Expiration: "01/20",
			Cvv:        123,
		}
		serv.EXPECT().Save(gomock.Any(), dto.SaveCreditCardDTO{
			Metadata: dto.CreditCardMetadata{
				BankName: "bank name",
			},
			CardNumber: validCardNumber,
			Expiration: "01/20",
			CVV:        123,
		}).Times(1).Return(fmt.Errorf("internal error"))

		// Execute.
		got, err := handler.SaveCreditCard(ctx, arg)

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, e.Code())
		assert.Nil(t, got)
	})

	t.Run("success save case", func(t *testing.T) {
		// Prepare.
		arg := &pb.SaveCreditCardRequest{
			Metadata:   &pb.CreditCardMetadata{BankName: "bank name"},
			CardNumber: validCardNumber,
			Expiration: "01/22",
			Cvv:        321,
		}
		serv.EXPECT().Save(gomock.Any(), dto.SaveCreditCardDTO{
			Metadata: dto.CreditCardMetadata{
				BankName: "bank name",
			},
			CardNumber: validCardNumber,
			Expiration: "01/22",
			CVV:        321,
		}).Times(1).Return(nil)

		// Execute.
		got, err := handler.SaveCreditCard(ctx, arg)

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
	})
}
