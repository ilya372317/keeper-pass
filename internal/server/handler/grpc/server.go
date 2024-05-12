package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/ilya372317/pass-keeper/internal/server/config"
	"github.com/ilya372317/pass-keeper/pkg/logger"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc"
)

type Server struct {
	srv  *grpc.Server
	conf config.Config
}

func New(cnfg config.Config,
	unaryInterceptors []grpc.UnaryServerInterceptor,
	streamInterceptors []grpc.StreamServerInterceptor,
) *Server {
	return &Server{
		srv: grpc.NewServer(
			grpc.ChainUnaryInterceptor(unaryInterceptors...),
			grpc.ChainStreamInterceptor(streamInterceptors...),
		),
		conf: cnfg}
}

func (s *Server) RegisterHandler(service pb.PassServiceServer) {
	pb.RegisterPassServiceServer(s.srv, service)
}

// StartAndListen start grpc server and listen incoming requests.
// When ctx.Done() will close, grpc server will gracefully shut down.
func (s *Server) StartAndListen(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.conf.GRPC.Host)
	if err != nil {
		return fmt.Errorf("failed start listen tcp connection on host [%s]: %w", s.conf.GRPC.Host, err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Infof("grpc server listen connection on host: %s", s.conf.GRPC.Host)
		if err = s.srv.Serve(lis); err != nil {
			logger.Log.Errorf("failed correct stop grpc server: %v", err)
		}
	}()

	<-ctx.Done()
	s.srv.GracefulStop()
	logger.Log.Info("grpc server was gracefully shutdown")

	return nil
}
