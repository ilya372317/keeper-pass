package passcleint

import (
	"context"
	"fmt"

	pb "github.com/ilya372317/pass-keeper/proto"
)

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
