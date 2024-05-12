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
