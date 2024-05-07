package v1

import (
	"context"

	pb "github.com/ilya372317/pass-keeper/proto"
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
		e := checkErr("credit-card", in.Id, err)
		return nil, status.Errorf(e.Code(), e.String())
	}

	return &pb.ShowCreditCardResponse{CreditCard: &pb.CreditCard{
		Id:         int64(creditCard.ID),
		CardNumber: creditCard.CardNumber,
		Expiration: creditCard.Expiration,
		Cvv:        int32(creditCard.CVV),
		Metadata:   &pb.CreditCardMetadata{BankName: creditCard.Metadata.BankName},
	}}, nil
}
