package v1

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ShowLoginPass(ctx context.Context, in *pb.ShowLoginPassRequest) (*pb.ShowLoginPassResponse, error) {
	data, err := h.loginPassService.Show(ctx, int(in.Id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "data with id [%d] not found", in.Id)
		}
		if errors.Is(err, domain.ErrNotSupportedOperation) {
			return nil, status.Errorf(codes.InvalidArgument, "try another method for show data with id [%d]", in.Id)
		}
		if errors.Is(err, domain.ErrAccesDenied) {
			return nil, status.Errorf(codes.PermissionDenied, "you can`t view data with id [%d]", in.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed show login pass: %v", err)
	}

	return &pb.ShowLoginPassResponse{LoginPass: &pb.LoginPass{
		Id:       int64(data.ID),
		Login:    data.Login,
		Password: data.Password,
		Metadata: &pb.LoginPassMetadata{Url: data.Metadata.URL},
	}}, nil
}
