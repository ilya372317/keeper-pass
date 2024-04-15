package v1

import (
	pb "github.com/ilya372317/pass-keeper/proto"
)

type AuthService interface {
}

type Handler struct {
	pb.UnimplementedPassServiceServer
	authService AuthService
}

func New(authService AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}
