package v1

import (
	pb "github.com/ilya372317/pass-keeper/proto"
)

type Handler struct {
	pb.UnimplementedPassServiceServer
}

func New() *Handler {
	return &Handler{}
}
