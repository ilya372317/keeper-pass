package v1

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) Auth(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	loginDTO := dto.LoginDTO{
		Email:    in.Email,
		Password: in.Password,
	}

	if err := dto.ValidateDTOWithRequired(loginDTO); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data: %v", err)
	}

	token, err := h.authService.Login(ctx, loginDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user with email [%s] not found", loginDTO.Email)
		}
		if errors.Is(err, domain.ErrInvalidPassword) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid password given")
		}
		return nil, status.Errorf(codes.Internal, "failed auth user: %v", err)
	}

	return &pb.AuthResponse{AccessToken: token}, nil
}
