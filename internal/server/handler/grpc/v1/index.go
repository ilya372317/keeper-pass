package v1

import (
	"context"

	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) Index(ctx context.Context, _ *pb.IndexRequest) (*pb.IndexResponse, error) {
	records, err := h.generalDataService.Index(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed get data list: %v", err)
	}

	responseDataList := make([]*pb.IndexItem, 0, len(records))

	for _, r := range records {
		responseDataList = append(responseDataList, &pb.IndexItem{
			Id:   r.ID,
			Info: r.Info,
			Type: int32(r.Kind),
		})
	}

	return &pb.IndexResponse{
		Items: responseDataList,
	}, nil
}
