package v1

import (
	"context"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	pb "github.com/ilya372317/pass-keeper/proto"
)

type AuthService interface {
	Login(context.Context, dto.LoginDTO) (string, error)
	Register(context.Context, dto.RegisterDTO) (*domain.User, error)
}

type loginPassService interface {
	Save(ctx context.Context, d dto.SaveLoginPassDTO) error
	Update(ctx context.Context, d dto.UpdateLoginPassDTO) error
	Show(ctx context.Context, id int) (domain.LoginPassData, error)
}

type Handler struct {
	pb.UnimplementedPassServiceServer
	authService      AuthService
	loginPassService loginPassService
}

func New(authService AuthService, loginPassService loginPassService) *Handler {
	return &Handler{
		authService:      authService,
		loginPassService: loginPassService,
	}
}
