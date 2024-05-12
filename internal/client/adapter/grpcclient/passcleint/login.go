package passcleint

import (
	"context"
	"fmt"

	pb "github.com/ilya372317/pass-keeper/proto"
)

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
