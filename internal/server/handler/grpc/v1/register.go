package v1

import (
	"context"

	pb "github.com/ilya372317/pass-keeper/proto"
)

func (h *Handler) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{}, nil
}
