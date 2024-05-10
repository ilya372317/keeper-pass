package v1

import (
	"context"

	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := h.generalDataService.Delete(ctx, in.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed delete data with id [%d]: %v", in.Id, err)
	}

	return &pb.DeleteResponse{}, nil
}
