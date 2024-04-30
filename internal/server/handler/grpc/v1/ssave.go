package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) SaveSimpleData(
	ctx context.Context,
	in *pb.SaveSimpleDataRequest,
) (
	*pb.SaveSimpleDataResponse,
	error,
) {
	d := dto.SaveSimpleDataDTO{
		Payload:  in.Payload,
		Metadata: in.Metadata,
		Type:     domain.Kind(in.Type),
	}
	if err := dto.ValidateDTOWithRequired(&d); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "given request is invalid: %v", err.Error())
	}

	data, err := h.dataService.SaveSimpleData(ctx, d)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed save data: %v", err.Error())
	}

	return &pb.SaveSimpleDataResponse{
		Data: &pb.Data{
			Payload:  data.Payload,
			Metadata: data.Metadata,
			Id:       int64(data.ID),
		},
	}, nil
}
