package passcleint

import (
	"context"
	"fmt"

	pb "github.com/ilya372317/pass-keeper/proto"
)

func (pc *PassClient) SaveLogin(ctx context.Context, login string, password string, url string) error {
	req := &pb.SaveLoginPassRequest{
		Metadata: &pb.LoginPassMetadata{
			Url: url,
		},
		Login:    login,
		Password: password,
	}

	if _, err := pc.c.SaveLoginPass(ctx, req); err != nil {
		return fmt.Errorf("grpc save login request failed: %w", err)
	}

	return nil
}

func (pc *PassClient) SaveCard(ctx context.Context, number, exp string, code int, bankName string) error {
	req := &pb.SaveCreditCardRequest{
		Metadata:   &pb.CreditCardMetadata{BankName: bankName},
		CardNumber: number,
		Expiration: exp,
		Cvv:        int32(code),
	}
	if _, err := pc.c.SaveCreditCard(ctx, req); err != nil {
		return fmt.Errorf("grpc save card request failed: %w", err)
	}

	return nil
}