package v1

import (
	"context"

	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/status"
)

func (h *Handler) ShowText(ctx context.Context, in *pb.ShowTextRequest) (*pb.ShowTextResponse, error) {
	text, err := h.textService.Show(ctx, in.Id)
	if err != nil {
		e := checkErr("text", in.Id, err)
		return nil, status.Errorf(e.Code(), e.String())
	}

	return &pb.ShowTextResponse{Text: &pb.Text{
		Id:       text.ID,
		Metadata: &pb.TextMetadata{Info: text.Metadata.Info},
		Data:     text.Data,
	}}, nil
}
