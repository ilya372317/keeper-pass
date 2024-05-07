package v1

import (
	"context"

	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/status"
)

func (h *Handler) ShowBinary(ctx context.Context, in *pb.ShowBinaryRequest) (*pb.ShowBinaryResponse, error) {
	binary, err := h.binaryService.Show(ctx, in.Id)
	if err != nil {
		e := checkErr("binary", in.Id, err)
		return nil, status.Errorf(e.Code(), e.String())
	}

	return &pb.ShowBinaryResponse{Binary: &pb.Binary{
		Id:       binary.ID,
		Metadata: &pb.BinaryMetadata{Info: binary.Metadata.Info},
		Data:     binary.Data,
	}}, nil
}
