package passcleint

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (pc *PassClient) All(ctx context.Context) ([]domain.ShortData, error) {
	req := &pb.IndexRequest{}
	resp, err := pc.c.Index(ctx, req)
	if err != nil {
		e, ok := status.FromError(err)
		if !ok {
			return nil, fmt.Errorf("invalid error recived from server: %w", err)
		}
		if e.Code() == codes.Unauthenticated {
			return nil, domain.ErrUnauthenticated
		}

		return nil, fmt.Errorf("failed get index data from server: %w", err)
	}

	shortDataRecords := make([]domain.ShortData, 0, len(resp.Items))

	for _, item := range resp.Items {
		shortDataRecords = append(shortDataRecords, domain.ShortData{
			ID:   item.Id,
			Info: item.Info,
			Kind: domain.Kind(item.Type),
		})
	}

	return shortDataRecords, nil
}
