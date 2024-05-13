package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) SaveCreditCard(
	ctx context.Context,
	in *pb.SaveCreditCardRequest,
) (
	*pb.SaveCreditCardResponse,
	error,
) {
	d := dto.SaveCreditCardDTO{
		CardNumber: in.CardNumber,
		Expiration: in.Expiration,
		CVV:        int(in.Cvv),
	}

	if in.Metadata != nil {
		d.Metadata.BankName = in.Metadata.BankName
	}

	if err := dto.ValidateDTOWithRequired(&d); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid credit card info given: %v", err)
	}

	if err := h.creditCardService.Save(ctx, d); err != nil {
		return nil, status.Errorf(codes.Internal, "failed save credit card data: %v", err)
	}

	return &pb.SaveCreditCardResponse{}, nil
}
