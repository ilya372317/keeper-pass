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

func TestHandler_ShowCreditCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := v1_mock.NewMockcreditCardService(ctrl)
	handler := Handler{
		creditCardService: serv,
	}
	ctx := context.Background()
	var arg int64 = 1

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Show(gomock.Any(), gomock.Any()).Times(1).Return(domain.CreditCardData{
			Metadata: domain.CreditCardMetadata{
				BankName: "string",
			},
			CardNumber: "123",
			Expiration: "123",
			CVV:        123,
			ID:         1,
		}, nil)

		// Execute.
		got, err := handler.ShowCreditCard(ctx, &pb.ShowCreditCardRequest{Id: arg})

		// Assert.
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, int64(1), got.CreditCard.Id)
		assert.Equal(t, "123", got.CreditCard.CardNumber)
		assert.Equal(t, "123", got.CreditCard.Expiration)
		assert.Equal(t, int32(123), got.CreditCard.Cvv)
		require.NotNil(t, got.CreditCard.Metadata)
		assert.Equal(t, "string", got.CreditCard.Metadata.BankName)
	})

	t.Run("data not found", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().
			Show(gomock.Any(), gomock.Any()).
			Times(1).
			Return(domain.CreditCardData{}, sql.ErrNoRows)

		// Execute.
		got, err := handler.ShowCreditCard(ctx, &pb.ShowCreditCardRequest{Id: arg})

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Nil(t, got)
		assert.Equal(t, codes.NotFound, e.Code())
	})

	t.Run("operation not supported", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().
			Show(gomock.Any(), gomock.Any()).
			Times(1).
			Return(domain.CreditCardData{}, domain.ErrNotSupportedOperation)

		// Execute.
		got, err := handler.ShowCreditCard(ctx, &pb.ShowCreditCardRequest{Id: arg})

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, e.Code())
		assert.Nil(t, got)
	})

	t.Run("internal error", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().
			Show(gomock.Any(), gomock.Any()).
			Times(1).
			Return(domain.CreditCardData{}, fmt.Errorf("internal"))

		// Execute.
		got, err := handler.ShowCreditCard(ctx, &pb.ShowCreditCardRequest{Id: arg})

		// Assert.
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, e.Code())
		assert.Nil(t, got)
	})
}
