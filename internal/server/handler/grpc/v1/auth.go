package v1

import (
	"context"

	pb "github.com/ilya372317/pass-keeper/proto"
)

func (h *Handler) Auth(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{AccessToken: "123"}, nil
}
