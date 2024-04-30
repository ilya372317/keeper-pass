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

type dataService interface {
	SaveSimpleData(context.Context, dto.SaveSimpleDataDTO) (*domain.Data, error)
}

type Handler struct {
	pb.UnimplementedPassServiceServer
	authService AuthService
	dataService dataService
}

func New(authService AuthService, dataService dataService) *Handler {
	return &Handler{
		authService: authService,
		dataService: dataService,
	}
}
