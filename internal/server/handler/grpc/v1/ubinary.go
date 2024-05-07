package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateBinary(ctx context.Context, in *pb.UpdateBinaryRequest) (*pb.UpdateBinaryResponse, error) {
	d := dto.UpdateBinaryDTO{
		Data: &in.Data,
		ID:   in.Id,
	}
	if in.Metadata != nil {
		d.Metadata = &dto.BinaryMetadata{Info: in.Metadata.Info}
	}

	if err := h.binaryService.Update(ctx, d); err != nil {
		return nil, status.Errorf(codes.Internal, "failed update binary data: %v", err)
	}

	return &pb.UpdateBinaryResponse{}, nil
}
