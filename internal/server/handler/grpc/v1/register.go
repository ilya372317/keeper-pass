package v1

import (
	"context"
	"errors"
	"time"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	registerDTO := dto.RegisterDTO{
		Email:    in.Email,
		Password: in.Password,
	}

	if err := dto.ValidateDTOWithRequired(&registerDTO); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument for registration given: %v", err)
	}

	user, err := h.authService.Register(ctx, registerDTO)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already registered: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed register user: %v", err)
	}

	response := pb.RegisterResponse{
		User: &pb.User{
			Id:        int64(user.ID),
			Email:     user.Email,
			CreatedAt: user.CreatedAT.Format(time.DateTime),
			UpdatedAt: user.UpdatedAT.Format(time.DateTime),
		},
	}

	return &response, nil
}
