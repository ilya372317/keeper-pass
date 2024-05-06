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

func (h *Handler) ShowText(ctx context.Context, in *pb.ShowTextRequest) (*pb.ShowTextResponse, error) {
	text, err := h.textService.Show(ctx, in.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "data with id [%d] not exists", in.Id)
		}
		if errors.Is(err, domain.ErrAccesDenied) {
			return nil, status.Errorf(codes.PermissionDenied, "you can`t view text data with id [%d]", in.Id)
		}
		if errors.Is(err, domain.ErrNotSupportedOperation) {
			return nil, status.Errorf(codes.InvalidArgument, "data with id [%d] is not text type", in.Id)
		}

		return nil, status.Errorf(codes.Internal, "failed get text data fro view: %v", err)
	}

	return &pb.ShowTextResponse{Text: &pb.Text{
		Id:       text.ID,
		Metadata: &pb.TextMetadata{Info: text.Metadata.Info},
		Data:     text.Data,
	}}, nil
}
