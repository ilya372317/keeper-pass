package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateCreditCard(
	ctx context.Context,
	in *pb.UpdateCreditCardRequest,
) (
	*pb.UpdateCreditCardResponse,
	error,
) {
	d := dto.UpdateCreditCardDTO{
		CardNumber: in.CardNumber,
		Expiration: in.Expiration,
		CVV:        in.Cvv,
		ID:         int(in.Id),
	}

	if in.Metadata != nil {
		d.Metadata = &dto.CreditCardMetadata{BankName: in.Metadata.BankName}
	}

	if err := dto.ValidateDTOWithRequired(&d); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid credit card data for update given: %v", err)
	}

	if err := h.creditCardService.Update(ctx, d); err != nil {
		e := checkErr("credit-card", in.Id, err)
		return nil, status.Errorf(e.Code(), e.String())
	}

	return &pb.UpdateCreditCardResponse{}, nil
}
