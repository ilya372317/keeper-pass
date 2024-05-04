package v1

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ShowCreditCard(
	ctx context.Context,
	in *pb.ShowCreditCardRequest,
) (
	*pb.ShowCreditCardResponse,
	error,
) {
	creditCard, err := h.creditCardService.Show(ctx, in.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "credit card with id [%d] not found", in.Id)
		}
		if errors.Is(err, domain.ErrNotSupportedOperation) {
			return nil, status.Errorf(codes.InvalidArgument, "data with id [%d] is not credit card", in.Id)
		}
		if errors.Is(err, domain.ErrAccesDenied) {
			return nil, status.Errorf(codes.PermissionDenied, "you can`t view this data")
		}

		return nil, status.Errorf(codes.Internal, "failed show credit card info: %v", err)
	}

	return &pb.ShowCreditCardResponse{CreditCard: &pb.CreditCard{
		Id:         int64(creditCard.ID),
		CardNumber: creditCard.CardNumber,
		Expiration: creditCard.Expiration,
		Cvv:        int32(creditCard.CVV),
		Metadata:   &pb.CreditCardMetadata{BankName: creditCard.Metadata.BankName},
	}}, nil
}
