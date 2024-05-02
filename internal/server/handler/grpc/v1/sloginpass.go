package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) SaveLoginPass(ctx context.Context, in *pb.SaveLoginPassRequest) (*pb.SaveLoginPassResponse, error) {
	d := dto.SaveLoginPassDTO{
		Login:    in.Login,
		Password: in.Password,
	}

	if in.Metadata != nil {
		d.Metadata = dto.LoginPassMetadata{URL: in.Metadata.Url}
	}

	if err := dto.ValidateDTOWithRequired(&d); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument given: %v", err)
	}

	if err := h.loginPassService.Save(ctx, d); err != nil {
		return nil, status.Errorf(codes.Internal, "failed save login pass data: %v", err)
	}

	return &pb.SaveLoginPassResponse{}, nil
}
