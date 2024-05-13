package passcleint

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PassClient struct {
	c pb.PassServiceClient
}

func New(client pb.PassServiceClient) *PassClient {
	return &PassClient{c: client}
}

func (pc *PassClient) All(ctx context.Context) ([]domain.ShortData, error) {
	req := &pb.IndexRequest{}
	resp, err := pc.c.Index(ctx, req)
	if err != nil {
		e, ok := status.FromError(err)
		if !ok {
			return nil, fmt.Errorf("invalid error recived from server: %w", err)
		}
		if e.Code() == codes.Unauthenticated {
			return nil, domain.ErrUnauthenticated
		}

		return nil, fmt.Errorf("failed get index data from server: %w", err)
	}

	shortDataRecords := make([]domain.ShortData, 0, len(resp.Items))

	for _, item := range resp.Items {
		shortDataRecords = append(shortDataRecords, domain.ShortData{
			ID:   item.Id,
			Info: item.Info,
			Kind: domain.Kind(item.Type),
		})
	}

	return shortDataRecords, nil
}

func (pc *PassClient) Login(ctx context.Context, email, password string) (string, error) {
	req := &pb.AuthRequest{
		Email:    email,
		Password: password,
	}
	resp, err := pc.c.Auth(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed get access token from grpc server: %w", err)
	}

	return resp.GetAccessToken(), nil
}

func (pc *PassClient) Register(ctx context.Context, email string, password string) error {
	req := &pb.RegisterRequest{
		Email:    email,
		Password: password,
	}
	if _, err := pc.c.Register(ctx, req); err != nil {
		return fmt.Errorf("error on server on register request: %w", err)
	}

	return nil
}

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

func (pc *PassClient) SaveText(ctx context.Context, info string, data string) error {
	req := &pb.SaveTextRequest{
		Metadata: &pb.TextMetadata{Info: info},
		Data:     data,
	}
	if _, err := pc.c.SaveText(ctx, req); err != nil {
		return fmt.Errorf("grpc save text request failed: %w", err)
	}

	return nil
}

func (pc *PassClient) SaveBinary(ctx context.Context, info string, data []byte) error {
	req := &pb.SaveBinaryRequest{
		Metadata: &pb.BinaryMetadata{Info: info},
		Data:     data,
	}

	if _, err := pc.c.SaveBinary(ctx, req); err != nil {
		return fmt.Errorf("grpc save binary request failed: %w", err)
	}

	return nil
}

func (pc *PassClient) Delete(ctx context.Context, id int) error {
	req := &pb.DeleteRequest{Id: int64(id)}
	if _, err := pc.c.Delete(ctx, req); err != nil {
		return fmt.Errorf("grpc delete failed: %w", err)
	}

	return nil
}

func (pc *PassClient) ShowLoginPass(ctx context.Context, id int) (domain.LoginPass, error) {
	req := &pb.ShowLoginPassRequest{Id: int64(id)}
	resp, err := pc.c.ShowLoginPass(ctx, req)
	if err != nil {
		return domain.LoginPass{}, fmt.Errorf("grpc show login pass request failed: %w", err)
	}
	lp := resp.LoginPass
	if lp == nil || lp.Metadata == nil {
		return domain.LoginPass{},
			fmt.Errorf("invalid show login pass response recived from grpc server: %w", err)
	}

	return domain.LoginPass{
		URL:      lp.Metadata.Url,
		Login:    lp.Login,
		Password: lp.Password,
		ID:       int(lp.Id),
	}, nil
}

func (pc *PassClient) ShowCard(ctx context.Context, id int) (domain.CreditCard, error) {
	req := &pb.ShowCreditCardRequest{Id: int64(id)}
	resp, err := pc.c.ShowCreditCard(ctx, req)
	if err != nil {
		return domain.CreditCard{}, fmt.Errorf("grpc show credit card request failed: %w", err)
	}

	cc := resp.CreditCard

	if cc == nil || cc.Metadata == nil {
		return domain.CreditCard{},
			fmt.Errorf("invalid show credit card response recived from grpc server: %w", err)
	}

	return domain.CreditCard{
		BankName:   cc.Metadata.BankName,
		CardNumber: cc.CardNumber,
		Exp:        cc.Expiration,
		Code:       int(cc.Cvv),
		ID:         int(cc.Id),
	}, nil
}

func (pc *PassClient) ShowText(ctx context.Context, id int) (domain.Text, error) {
	req := &pb.ShowTextRequest{Id: int64(id)}
	resp, err := pc.c.ShowText(ctx, req)
	if err != nil {
		return domain.Text{}, fmt.Errorf("grpc show text request failed: %w", err)
	}

	t := resp.Text
	if t == nil || t.Metadata == nil {
		return domain.Text{}, fmt.Errorf("invalid show text response recived from grpc server: %w", err)
	}

	return domain.Text{
		Info: t.Metadata.Info,
		Data: t.Data,
		ID:   int(t.Id),
	}, nil
}

func (pc *PassClient) ShowBinary(ctx context.Context, id int) (domain.Binary, error) {
	req := &pb.ShowBinaryRequest{Id: int64(id)}
	resp, err := pc.c.ShowBinary(ctx, req)
	if err != nil {
		return domain.Binary{}, fmt.Errorf("grpc show binary request failed: %w", err)
	}

	b := resp.Binary

	if b == nil || b.Metadata == nil {
		return domain.Binary{}, fmt.Errorf("invalid show binary response from grpc server recived: %w", err)
	}

	return domain.Binary{
		Info: b.Metadata.Info,
		Data: b.Data,
		ID:   int(b.Id),
	}, nil
}
