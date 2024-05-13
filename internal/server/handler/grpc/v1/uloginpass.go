package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateLoginPass(
	ctx context.Context,
	in *pb.UpdateLoginPassRequest,
) (
	*pb.UpdateLoginPassResponse,
	error,
) {
	d := dto.UpdateLoginPassDTO{
		ID:       in.Id,
		Login:    in.Login,
		Password: in.Password,
	}
	if in.Metadata != nil {
		d.Metadata = &dto.LoginPassMetadata{URL: in.Metadata.Url}
	}

	if err := dto.ValidateDTOWithRequired(&d); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument given for update login pass: %v", err)
	}

	if err := h.loginPassService.Update(ctx, d); err != nil {
		e := checkErr("loginpass", in.Id, err)
		return nil, status.Errorf(e.Code(), e.String())
	}

	return &pb.UpdateLoginPassResponse{}, nil
}
