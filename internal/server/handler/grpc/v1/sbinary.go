package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) SaveBinary(ctx context.Context, in *pb.SaveBinaryRequest) (*pb.SaveBinaryResponse, error) {
	d := dto.SaveBinaryDTO{
		Data: in.Data,
	}
	if in.Metadata != nil {
		d.Metadata.Info = in.Metadata.Info
	}

	if err := h.binaryService.Save(ctx, d); err != nil {
		return nil, status.Errorf(codes.Internal, "failed save data: %v", err)
	}

	return &pb.SaveBinaryResponse{}, nil
}
