package v1

import (
	"context"
	"errors"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
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
		if errors.Is(err, domain.ErrAccesDenied) {
			return nil, status.Errorf(codes.PermissionDenied, "you can`t edit data belongs to other person")
		}

		return nil, status.Errorf(codes.Internal, "failed update data with id [%d]: %v", in.Id, err)
	}

	return &pb.UpdateCreditCardResponse{}, nil
}
