package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/ilya372317/pass-keeper/internal/server/config"
	v1 "github.com/ilya372317/pass-keeper/internal/server/handler/grpc/v1"
	"github.com/ilya372317/pass-keeper/internal/server/logger"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc"
)

type Server struct {
	srv  *grpc.Server
	conf config.Config
}

func New(cnfg config.Config) *Server {
	return &Server{
		srv: grpc.NewServer(
			grpc.ChainUnaryInterceptor(getUnaryInterceptors()...),
			grpc.ChainStreamInterceptor(getStreamInterceptors()...),
		),
		conf: cnfg}
}

// StartAndListen start grpc server and listen incoming requests.
// When ctx.Done() will close, grpc server will gracefully shut down.
func (s *Server) StartAndListen(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.conf.GRPC.Host)
	if err != nil {
		return fmt.Errorf("failed start listen tcp connection on host [%s]: %w", s.conf.GRPC.Host, err)
	}

	pb.RegisterPassServiceServer(s.srv, v1.New())
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

func getUnaryInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(logger.InterceptorLogger()),
	}
}

func getStreamInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{}
}
