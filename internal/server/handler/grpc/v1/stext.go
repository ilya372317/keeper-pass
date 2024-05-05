package v1

import (
	"context"

	pb "github.com/ilya372317/pass-keeper/proto"
)

func (h *Handler) SaveText(ctx context.Context, in *pb.SaveTextRequest) (*pb.SaveTextResponse, error) {
	return &pb.SaveTextResponse{}, nil
}
