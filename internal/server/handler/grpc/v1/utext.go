package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
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
		e := checkErr("text", in.Id, err)
		return nil, status.Errorf(e.Code(), e.String())
	}

	return &pb.UpdateTextResponse{}, nil
}
