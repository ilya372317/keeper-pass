package grpc

import "google.golang.org/grpc"

type GRPCServer struct {
	srv *grpc.Server
}

func New() *GRPCServer {
	return &GRPCServer{srv: grpc.NewServer()}
}

func (g *GRPCServer) Start() {
}
