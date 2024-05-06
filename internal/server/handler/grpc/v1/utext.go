package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateText(ctx context.Context, in *pb.UpdateTextRequest) (*pb.UpdateTextResponse, error) {
	d := dto.UpdateTextDTO{
		Data: in.Data,
		ID:   in.Id,
	}
	if in.Metadata != nil {
		d.Metadata = &dto.TextMetadata{Info: in.Metadata.Info}
	}

	if err := h.textService.Update(ctx, d); err != nil {
		return nil, status.Errorf(codes.Internal, "failed update text: %v", err)
	}

	return &pb.UpdateTextResponse{}, nil
}
